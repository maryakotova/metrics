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
	var interval int
	for range agent.pollTicker.C {
		interval += int(agent.PollInterval)
		if interval == int(agent.ReportInterval) {
			interval = 0
			metrics := agent.collectRuntimeMetrics()
			agent.sendQueue <- Metrics{Metrics: metrics}
		}
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
	// agent.wg.Add(1)
	// defer agent.wg.Done()
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
		result, err := <-agent.resultQueue
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		// канал закрыт
		if err {
			return
		}
	}
}
