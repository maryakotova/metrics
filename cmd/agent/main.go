package main

import (
	"fmt"
	"time"

	"github.com/maryakotova/metrics/internal/collector"
	"github.com/maryakotova/metrics/internal/sender"
)

var (
	serverAddress  string
	pollInterval   int64
	reportInterval int64
)

func main() {

	parseFlags()

	n := int64(0)
	collector.SetPollCountInitial()

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		n += pollInterval

		metrics := collector.CollectMetrics()

		if reportInterval == n {
			n = 0
			collector.SetPollCountInitial()

			for metricName, metricValue := range metrics {
				var metricType string
				if metricName == "PollCount" {
					metricType = "counter"
				} else {
					metricType = "gauge"
				}
				err := sender.SendMetric(serverAddress, metricType, metricName, metricValue)
				if err != nil {
					fmt.Printf("Ошибка при отправке метрики %s: %s\n", metricName, err)
				}
			}
		}
	}
}
