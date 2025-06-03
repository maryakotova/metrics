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
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"metrics/internal/config"
	"metrics/internal/constants"
	"metrics/internal/controller"
	"metrics/internal/decryptmiddleware"
	"metrics/internal/filetransfer"
	"metrics/internal/handlers"
	"metrics/internal/logger"
	"metrics/internal/middleware"
	"metrics/internal/storage"
	"metrics/internal/trustedip"
	"metrics/internal/worker"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// глобальные переменные с информацией о версии
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {

	printVersionInfo()

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

	handler := handlers.NewServer(config, writer, log, controller)

	if !config.IsSyncStore() && config.IsStoreInFileEnabled() {
		ticker := time.NewTicker(time.Duration(config.Server.StoreInterval) * time.Second)
		defer ticker.Stop()
		task := func() {
			metrics := storage.GetAllMetricsInJSON()
			writer.WriteMetrics(metrics...)
		}
		worker.TriggerGoFunc(ticker, task)
	}

	dmw := decryptmiddleware.NewDecrypteMW(config, log)
	ipmw := trustedip.NewCheckIPMW(config, log)

	router := createRouter(handler, dmw, ipmw)

	go func() {
		log.Info("pprof listening on :6060")
		http.ListenAndServe("localhost:6060", nil) // <- pprof listening on :6060
	}()

	var server = http.Server{
		Addr:    config.Server.ServerAddress,
		Handler: router,
	}

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})

	// канал для перенаправления прерываний
	sigint := make(chan os.Signal, 1)

	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		// читаем из канала прерываний
		<-sigint

		// запускаем процедуру graceful shutdown
		if err := server.Shutdown(context.Background()); err != nil {
			err = fmt.Errorf("HTTP server Shutdown: %w", err)
			log.Error(err.Error())
		}
		// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		err = fmt.Errorf("HTTP server ListenAndServe: %w", err)
		log.Fatal(err.Error())
	}

	<-idleConnsClosed

	// выделяем время на сохранение данных либо сохраняем
	if !config.IsSyncStore() && config.IsStoreInFileEnabled() {
		metrics := storage.GetAllMetricsInJSON()
		writer.WriteMetrics(metrics...)
	} else {
		time.After(constants.ShutdownTimeout)
	}

}

func printVersionInfo() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}

func createRouter(handler *handlers.Server, dmw *decryptmiddleware.DecryptMiddleware, ipmw *trustedip.CheckIPMiddleware) *chi.Mux {
	router := chi.NewRouter()
	router.Use()

	router.Get("/", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandleGetAllMetrics)))))

	router.Get("/value/{metricType}/{metricName}", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandleGetOneMetric)))))
	router.Post("/value/", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandleGetOneMetricViaJSON)))))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandleMetricUpdate)))))
	router.Post("/update/", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandleMetricUpdateViaJSON)))))
	router.Post("/updates/", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandleMetricUpdates)))))

	router.Get("/ping", logger.WithLogging(ipmw.CheckIP(middleware.GzipMiddleware(dmw.Decrypte(handler.HandlePing)))))

	return router
}
