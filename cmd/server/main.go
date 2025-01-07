package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/maryakotova/metrics/internal/handlers"
	"github.com/maryakotova/metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	netAddress := flag.String("a", "localhost:8080", "Адрес и порт для HTTP-сервера")
	flag.Parse()

	memStorage := storage.NewMemStorage()
	server := handlers.NewServer(memStorage)

	router := chi.NewRouter()
	router.Use()

	router.Get("/", server.HandleGetAllMetrics)
	router.Get("/value/{metricType}/{metricName}", server.HandleGetOneMetric)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", server.HandleMetricUpdate)

	err := http.ListenAndServe(*netAddress, router)
	if err != nil {
		panic(err)
	}

	fmt.Println("Running server on ", netAddress)

}
