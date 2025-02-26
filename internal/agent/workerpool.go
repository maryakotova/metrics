package agent

import "fmt"

type Metrics struct {
	Metrics map[string]interface{}
}

type Result struct {
	Metrics map[string]interface{}
	Error   error
}

func (agent *Agent) CollectRuntimeMetricsAtInterval() {
	for range agent.pollTicker.C {
		metrics := agent.collectRuntimeMetrics()
		agent.sendQueue <- Metrics{Metrics: metrics}
	}
}

func (agent *Agent) CollectAdditionalMetricsAtInterval() {
	for range agent.pollTicker.C {
		metrics, err := agent.collectAdditionalMetrics()
		if err != nil {
			//писать в лог
		} else {
			agent.sendQueue <- Metrics{Metrics: metrics}
		}
	}
}

func (agent *Agent) Worker() {
	for range agent.reportTicker.C {
		for metrics := range agent.sendQueue {
			metricsForSend := agent.PrepareMetrics(metrics.Metrics)
			err := agent.SendMetricsBatch(metricsForSend)
			if err != nil {
				agent.resultQueue <- Result{Metrics: metrics.Metrics, Error: err}
			}
		}
	}
}

func (agent *Agent) HandleErrors() {
	for {
		result, open := <-agent.resultQueue
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		// канал закрыт
		if !open {
			return
		}
	}
}
