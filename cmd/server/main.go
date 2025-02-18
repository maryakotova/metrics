package main

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/controller"
	"github.com/maryakotova/metrics/internal/database/postgres"
	"github.com/maryakotova/metrics/internal/filetransfer"
	"github.com/maryakotova/metrics/internal/handlers"
	"github.com/maryakotova/metrics/internal/inmemory"
	"github.com/maryakotova/metrics/internal/logger"
	"github.com/maryakotova/metrics/internal/middleware"
	"github.com/maryakotova/metrics/internal/worker"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	parseFlags()

	log, err := logger.Initialize("")
	if err != nil {
		panic(err)
	}

	var ctrl *controller.Controller

	var server *handlers.Server

	if dbDsn != "" {
		db, err := sql.Open("pgx", dbDsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		postgresStorage := postgres.NewPostgresStorage(db, log, constants.RetryCount)
		err = postgresStorage.Bootstrap(context.TODO())
		if err != nil {
			panic(err)
		}

		ctrl = controller.NewController(postgresStorage, log)
		server = handlers.NewServer(false, nil, log, ctrl)
	}

	if filePath != "" {
		memStorage := inmemory.NewMemStorage()
		if restore {
			memStorage.UploadData(filePath)
		}
		writer, err := filetransfer.NewFileWriter(filePath)
		if err != nil {
			panic(err)
		}
		defer writer.Close()

		var syncFileWrite bool
		if interval == 0 {
			syncFileWrite = true
		} else {
			ticker := time.NewTicker(time.Duration(interval) * time.Second)
			defer ticker.Stop()
			task := func() {
				metrics := memStorage.GetAllMetricsInJSON()
				writer.WriteMetrics(metrics...)
			}
			worker.TriggerGoFunc(ticker, task)
		}

		ctrl = controller.NewController(memStorage, log)
		server = handlers.NewServer(syncFileWrite, writer, log, ctrl)
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

	err = http.ListenAndServe(netAddress, router)
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
