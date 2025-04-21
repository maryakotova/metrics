package agent

import (
	"fmt"
	"log/slog"
)

type Metrics struct {
	Metrics map[string]interface{}
}

type Result struct {
	Metrics map[string]interface{}
	Error   error
}

func (agent *Agent) CollectRuntimeMetricsAtInterval() {
	for range agent.pollTicker.C {
		agent.collectRuntimeMetrics()
	}
}

func (agent *Agent) CollectAdditionalMetricsAtInterval() {
	for range agent.pollTicker.C {
		agent.collectAdditionalMetrics()

	}
}

func (agent *Agent) PublishMetrics() {
	for range agent.reportTicker.C {
		agent.mutex.Lock()
		defer agent.mutex.Unlock()
		agent.sendQueue <- Metrics{Metrics: agent.metrics}
		// slog.Info(fmt.Sprintf("metrics published: %v", agent.metrics))
		slog.Info("metrics published")
		agent.ResetMetricsStorage()
		agent.setPollCountInitial()

	}

}

func (agent *Agent) Worker(w int) {
	// for metrics := range agent.sendQueue {
	slog.Info(fmt.Sprintf("worker %v started", w))
	defer agent.WG.Done()
	for {
		metrics := <-agent.sendQueue
		// slog.Info(fmt.Sprintf("Worker %v: metrics read (%v)", w, metrics))
		slog.Info(fmt.Sprintf("Worker %v: metrics read", w))
		metricsForSend := agent.PrepareMetrics(metrics.Metrics)
		err := agent.SendMetricsBatch(metricsForSend)
		agent.resultQueue <- Result{Metrics: metrics.Metrics, Error: err}
	}
}

func (agent *Agent) HandleErrors() {
	defer agent.WG.Done()
	for result := range agent.resultQueue {
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		} else {
			fmt.Println("метрики отправлены успешно")
		}
	}
}
