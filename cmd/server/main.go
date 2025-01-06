package main

import (
	"net/http"

	"github.com/maryakotova/metrics/internal/handlers"
	"github.com/maryakotova/metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	memStorage := storage.NewMemStorage()
	server := handlers.NewServer(memStorage)

	router := chi.NewRouter()
	router.Use()

	router.Get("/", server.HandleGetAllMetrics)
	router.Get("/value/{metricType}/{metricName}", server.HandleGetOneMetric)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", server.HandleMetricUpdate)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
