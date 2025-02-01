package main

import (
	"net/http"
	"strings"

	"github.com/maryakotova/metrics/internal/filetransfer"
	"github.com/maryakotova/metrics/internal/handlers"
	"github.com/maryakotova/metrics/internal/logger"
	"github.com/maryakotova/metrics/internal/middleware"
	"github.com/maryakotova/metrics/internal/storage"
	"github.com/maryakotova/metrics/internal/worker"

	"github.com/go-chi/chi/v5"
)

func main() {

	parseFlags()

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := logger.Initialize(""); err != nil {
		panic(err)
	}

	memStorage := storage.NewMemStorage()

	if restore {
		memStorage.UploadData(filePath)
	}

	writer, err := filetransfer.NewFileWriter(filePath)
	if err != nil {
		panic(err)
	}

	var syncFileWrite bool
	if interval == 0 {
		syncFileWrite = true
	} else {
		task := func() {
			metrics := memStorage.GetAllMetricsInJSON()
			writer.WriteMetrics(&metrics)
		}
		worker.InitPeriodicFunc(interval, task)
	}

	server := handlers.NewServer(memStorage, syncFileWrite, writer)

	router := chi.NewRouter()
	router.Use()

	router.Get("/", logger.WithLogging(gzipMiddleware(server.HandleGetAllMetrics)))

	router.Get("/value/{metricType}/{metricName}", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetric)))
	router.Post("/value/", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetricViaJSON)))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdate)))
	router.Post("/update/", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdateViaJSON)))

	err = http.ListenAndServe(netAddress, router)
	if err != nil {
		panic(err)
	}

	// <-signalChan

	// if interval != 0 {
	// 	metrics := memStorage.GetAllMetricsInJSON()
	// 	writer.WriteMetrics(&metrics)
	// }
	// writer.Close()
	// os.Exit(0)
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
