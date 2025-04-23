// Сервер для сбора рантайм-метрик собирает репорты от агентов по протоколу HTTP
//
// Сервер может принимать и хранить произвольные метрики двух типов gauge и counter.
//
// # Запуск сервера
//
//	В зависимости от флага -d и переменной окружения DATABASE_DSN метрики сохраняются в БД PostgreSQL или хранятся в памяти.
//	Флаг -i и переменная окружения STORE_INTERVAL отвечают за интервал времени в секундах, по истечении которого текущие показания сервера сохраняются в файл.
//	Флаг -f, переменная окружения FILE_STORAGE_PATH отвечают за полное имя файла, куда сохраняются текущие значения.
//	Флаг -r, переменная окружения RESTORE определяют загружать или нет ранее сохранённые значения из указанного файла при старте сервера.
//	Флаг -k и переменная окружения KEY содержат в себе секретный ключ для хэширования данных.
//	Флаг -d и переменная окружения DATABASE_DSN содержат адресом подключения к БД.
//
// # Эндпоинты
//
//	GET / - вывод всех метрик в
//	GET /value/{metricType}/{metricName} - возврат текущего значения метрики в текстовом виде
//	GET /ping - при запросе проверяет соединение с базой данных
//	POST /value/ - возврат текущего значения метрики в формате JSON
//	POST /update/{metricType}/{metricName}/{metricValue} - получение метрики с использованием Content-Type: text/plain
//	POST /update/ - получение метрики с использованием Content-Type: application/json
//	POST /updates/ - получение множества метрики с использованием Content-Type: application/json
package main

import (
	"net/http"
	"time"

	_ "net/http/pprof"

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

	router.Get("/", logger.WithLogging(middleware.GzipMiddleware(server.HandleGetAllMetrics)))

	router.Get("/value/{metricType}/{metricName}", logger.WithLogging(middleware.GzipMiddleware(server.HandleGetOneMetric)))
	router.Post("/value/", logger.WithLogging(middleware.GzipMiddleware(server.HandleGetOneMetricViaJSON)))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(middleware.GzipMiddleware(server.HandleMetricUpdate)))
	router.Post("/update/", logger.WithLogging(middleware.GzipMiddleware(server.HandleMetricUpdateViaJSON)))
	router.Post("/updates/", logger.WithLogging(middleware.GzipMiddleware(server.HandleMetricUpdates)))

	router.Get("/ping", logger.WithLogging(middleware.GzipMiddleware(server.HandlePing)))

	go func() {
		log.Info("pprof listening on :6060")
		http.ListenAndServe("localhost:6060", nil) // <- pprof listening on :6060
	}()

	err = http.ListenAndServe(config.Server.ServerAddress, router)
	if err != nil {
		panic(err)
	}

}

// func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ow := w

// 		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
// 		supportsGzipJSON := strings.Contains(r.Header.Get("Accept"), "application/json")
// 		supportsGzipHTML := strings.Contains(r.Header.Get("Accept"), "text/html")
// 		if supportsGzip && (supportsGzipJSON || supportsGzipHTML) {
// 			cw := middleware.NewCompressWriter(w)
// 			ow = cw
// 			defer cw.Close()
// 			ow.Header().Set("Content-Encoding", "gzip")
// 		}

// 		sendsGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
// 		if sendsGzip {
// 			cr, err := middleware.NewCompressReader(r.Body)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 			r.Body = cr
// 			defer cr.Close()
// 		}

// 		h.ServeHTTP(ow, r)
// 	}
// }
