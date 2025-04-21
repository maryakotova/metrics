package main

import (
	"net/http"
	"strings"
	"time"

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

	// var ctrl *controller.Controller

	// var server *handlers.Server

	// if dbDsn != "" {
	// 	db, err := sql.Open("pgx", dbDsn)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer db.Close()

	// 	postgresStorage := postgres.NewPostgresStorage(db, log, constants.RetryCount)
	// 	err = postgresStorage.Bootstrap(context.TODO())
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	ctrl = controller.NewController(postgresStorage, log)
	// 	server = handlers.NewServer(false, nil, log, ctrl)
	// }

	// if filePath != "" {
	// 	memStorage := inmemory.NewMemStorage()
	// 	if restore {
	// 		memStorage.UploadData(filePath)
	// 	}
	// 	writer, err := filetransfer.NewFileWriter(filePath)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer writer.Close()

	// 	var syncFileWrite bool
	// 	if interval == 0 {
	// 		syncFileWrite = true
	// 	} else {
	// 		ticker := time.NewTicker(time.Duration(interval) * time.Second)
	// 		defer ticker.Stop()
	// 		task := func() {
	// 			metrics := memStorage.GetAllMetricsInJSON()
	// 			writer.WriteMetrics(metrics...)
	// 		}
	// 		worker.TriggerGoFunc(ticker, task)
	// 	}

	// 	ctrl = controller.NewController(memStorage, log)
	// 	server = handlers.NewServer(syncFileWrite, writer, log, ctrl)
	// }

	router := chi.NewRouter()
	router.Use()

	router.Get("/", logger.WithLogging(gzipMiddleware(server.HandleGetAllMetrics)))

	router.Get("/value/{metricType}/{metricName}", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetric)))
	router.Post("/value/", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetricViaJSON)))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdate)))
	router.Post("/update/", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdateViaJSON)))
	router.Post("/updates/", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdates)))

	router.Get("/ping", logger.WithLogging(gzipMiddleware(server.HandlePing)))

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
