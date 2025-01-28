package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/maryakotova/metrics/internal/filetransfer"
	"github.com/maryakotova/metrics/internal/handlers"
	"github.com/maryakotova/metrics/internal/logger"
	"github.com/maryakotova/metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	parseFlags()

	if err := logger.Initialize(""); err != nil {
		panic(err)
	}

	writer, err := filetransfer.NewFileWriter(filePath)
	if err != nil {
		panic(err)
	}

	memStorage := storage.NewMemStorage()

	var syncFileWrite bool
	if interval == 0 {
		syncFileWrite = true
	} else {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				metrics := memStorage.GetAllMetricsInJSON()
				writer.WriteMetrics(&metrics)
			}

		}()
	}

	server := handlers.NewServer(memStorage, syncFileWrite, writer)

	UploadData(memStorage)

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

	defer writer.Close()
	if interval != 0 {
		defer func() {
			metrics := memStorage.GetAllMetricsInJSON()
			writer.WriteMetrics(&metrics)
		}()
	}
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		supportsGzipJSON := strings.Contains(r.Header.Get("Accept"), "application/json")
		supportsGzipHTML := strings.Contains(r.Header.Get("Accept"), "text/html")
		if supportsGzip && (supportsGzipJSON || supportsGzipHTML) {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
			ow.Header().Set("Content-Encoding", "gzip")
		}

		sendsGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
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
