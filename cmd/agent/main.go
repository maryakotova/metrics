package main

import (
	"github.com/maryakotova/metrics/internal/agent"
)

func main() {

	cfg, err := agent.ParseFlags()
	if err != nil {
		panic(err)
	}

	agent := agent.New(cfg)

	// один раз в PollInterval секунд в очередь добавляются данные
	// правильно? или нужно раз в PollInterval секунд сохранять их в локальное хранилище и раз в ReportInterval секунд добавлять в очередь?
	go agent.CollectRuntimeMetricsAtInterval()
	go agent.CollectAdditionalMetricsAtInterval()

	for w := 0; w <= int(agent.RateLimit); w++ {
		go agent.Worker()
	}

	agent.HandleErrors()

}
