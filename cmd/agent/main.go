package main

import (
	"metrics/internal/agent"
)

func main() {

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

	agent.WG.Wait()

}
