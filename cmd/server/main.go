package main

import (
	"net/http"
	"strings"

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

	memStorage := storage.NewMemStorage()
	server := handlers.NewServer(memStorage)

	router := chi.NewRouter()
	router.Use()

	router.Get("/", logger.WithLogging(server.HandleGetAllMetrics))

	router.Get("/value/{metricType}/{metricName}", logger.WithLogging(server.HandleGetOneMetric))
	router.Post("/value/", logger.WithLogging(gzipMiddleware(server.HandleGetOneMetricViaJSON)))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdate)))
	router.Post("/update/", logger.WithLogging(gzipMiddleware(server.HandleMetricUpdateViaJSON)))

	err := http.ListenAndServe(netAddress, router)
	if err != nil {
		panic(err)
	}
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		// contTypeJSON := strings.Contains(ow.Header().Get("Content-Type"), "application/json")
		// contTypeHTML := strings.Contains(ow.Header().Get("Content-Type"), "text/html")

		if supportsGzip { //&& (contTypeJSON || contTypeHTML) {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
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
