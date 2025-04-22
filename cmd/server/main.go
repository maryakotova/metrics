package main

import (
	"net/http"
	"strings"
	"time"

	//_ "net/http/pprof"

	"metrics/internal/config"
	"metrics/internal/controller"
	"metrics/internal/filetransfer"
	"metrics/internal/handlers"
	"metrics/internal/logger"
	"metrics/internal/middleware"
	"metrics/internal/storage"
	"metrics/internal/worker"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	flags, err := config.ParseFlags()
	if err != nil {
		panic(err)
	}

	log, err := logger.Initialize("")
	if err != nil {
		panic(err)
	}

	config := config.NewConfig(flags)

	factory := &storage.StorageFactory{}

	storage, err := factory.NewStorage(config, log)
	if err != nil {
		panic(err)
	}

	controller := controller.NewController(storage, log)

	var writer *filetransfer.FileWriter

	if config.IsStoreInFileEnabled() {
		writer, err = filetransfer.NewFileWriter(config.Server.FileStoragePath)
		if err != nil {
			panic(err)
		}
		defer writer.Close()
	}

	server := handlers.NewServer(config, writer, log, controller)

	if !config.IsSyncStore() && config.IsStoreInFileEnabled() {
		ticker := time.NewTicker(time.Duration(config.Server.StoreInterval) * time.Second)
		defer ticker.Stop()
		task := func() {
			metrics := storage.GetAllMetricsInJSON()
			writer.WriteMetrics(metrics...)
		}
		worker.TriggerGoFunc(ticker, task)
	}

	router := chi.NewRouter()
	router.Use()

	router.Get("/", logger.WithLogging(gzipMiddleware(server.HandleGetAllMetrics)))

	router.Get("/value/{metricType}/{metricName}", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetric)))
	router.Post("/value/", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetricViaJSON)))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdate)))
	router.Post("/update/", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdateViaJSON)))
	router.Post("/updates/", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdates)))

	router.Get("/ping", logger.WithLogging(gzipMiddleware(server.HandlePing)))

	go func() {
		log.Info("pprof listening on :6060")
		http.ListenAndServe("localhost:6060", nil) // <- DefaultServeMux
	}()

	err = http.ListenAndServe(config.Server.ServerAddress, router)
	if err != nil {
		panic(err)
	}

}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		supportsGzipJSON := strings.Contains(r.Header.Get("Accept"), "application/json")
		supportsGzipHTML := strings.Contains(r.Header.Get("Accept"), "text/html")
		if supportsGzip && (supportsGzipJSON || supportsGzipHTML) {
			cw := middleware.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
			ow.Header().Set("Content-Encoding", "gzip")
		}

		sendsGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if sendsGzip {
			cr, err := middleware.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}
}
