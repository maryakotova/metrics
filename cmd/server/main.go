package main

import (
	"net/http"

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
	router.Post("/value/", logger.WithLogging(server.HandleGetOneMetricViaJSON))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(server.HandleMetricUpdate))
	router.Post("/update/", logger.WithLogging(server.HandleMetricUpdateViaJSON))

	err := http.ListenAndServe(netAddress, router)
	if err != nil {
		panic(err)
	}
}
