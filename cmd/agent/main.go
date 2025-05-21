// Агент (HTTP-клиент) для сбора рантайм-метрик и их последующей отправки на сервер по протоколу HTTP.
//
// Агент собирает метрики двух типов: gauge и counter.
// В качестве источника метрик использованы пакеты runtime и gopsutil.
//
// Приложение обновляет метрики из пакета runtime с заданной во флаге -p или в переменной окружения  POLL_INTERVAL в частотой.
// Отправка метрик на сервер происходит с заданной во флаге -r или в переменной окружения  REPORT_INTERVAL частотой.
// Агент получает адрес эндпоинта HTTP-сервера из флага -a или переменной окружения ADDRESS.
// Флаг -k и переменная окружения KEY содержат в себе секретный ключ для хэширования данных.
// Количество одновременно исходящих запросов на сервер задается через флаг -l и переменную окружения RATE_LIMIT.
package main

import (
	"fmt"
	"metrics/internal/agent"

	"net/http"
	_ "net/http/pprof"
)

// глобальные переменные с информацией о версии
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {

	printVersionInfo()

	cfg, err := agent.ParseFlags()
	if err != nil {
		panic(err)
	}

	agent := agent.New(cfg)

	go agent.CollectRuntimeMetricsAtInterval()
	go agent.CollectAdditionalMetricsAtInterval()
	go agent.PublishMetrics()

	for w := range int(agent.RateLimit) {
		agent.WG.Add(1)
		go agent.Worker(w)
	}

	agent.WG.Add(1)
	go agent.HandleErrors()

	go func() {
		// log.Info("pprof listening on :6060")
		http.ListenAndServe("localhost:6061", nil) // <- DefaultServeMux
	}()

	agent.WG.Wait()

}

func printVersionInfo() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}
