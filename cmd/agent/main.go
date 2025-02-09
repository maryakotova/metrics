package main

import (
	"fmt"
	"time"

	"github.com/maryakotova/metrics/internal/collector"
	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/sender"
)

var (
	serverAddress  string
	pollInterval   int64
	reportInterval int64
)

func main() {

	parseFlags()

	// n := int64(0)
	// collector.SetPollCountInitial()

	// for {
	// 	time.Sleep(time.Duration(pollInterval) * time.Second)
	// 	n += pollInterval

	// 	metrics := collector.CollectMetrics()

	// 	if len(metrics) > 0 {
	// 		if reportInterval == n {
	// 			n = 0
	// 			collector.SetPollCountInitial()

	// 			if err := sender.SendMetrics(serverAddress, metrics); err != nil {
	// 				fmt.Printf("Ошибка при отправке метрик: %s\n", err)
	// 			}
	// 		}
	// 	}
	// }

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
				if metricName == constants.PollCount {
					metricType = constants.Counter
				} else {
					metricType = constants.Gauge
				}
				err := sender.SendMetric(serverAddress, metricType, metricName, metricValue)
				if err != nil {
					fmt.Printf("Ошибка при отправке метрики %s: %s\n", metricName, err)
				}
			}
		}
	}
}
