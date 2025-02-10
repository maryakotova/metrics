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

		metrics := collector.CollectMetricsForBatch()

		if len(metrics) > 0 {
			if reportInterval == n {
				n = 0
				collector.SetPollCountInitial()

				if err := sender.SendMetricsBatch(serverAddress, metrics); err != nil {
					fmt.Printf("Ошибка при отправке метрик: %s\n", err)
				}
			}
		}
	}

	// n := int64(0)
	// collector.SetPollCountInitial()

	// for {
	// 	time.Sleep(time.Duration(pollInterval) * time.Second)
	// 	n += pollInterval

	// 	metrics := collector.CollectMetrics()

	// 	if reportInterval == n {
	// 		n = 0
	// 		collector.SetPollCountInitial()

	// 		retries := 0
	// 		for retries < 4 {

	// 			err := sender.SendMetrics(serverAddress, metrics)
	// 			if err == nil {
	// 				break
	// 			}

	// 			var opErr *net.OpError
	// 			if !errors.As(err, &opErr) {
	// 				break
	// 			}

	// 			if retries == 3 {
	// 				fmt.Println("ошибка соединения: %w", err)
	// 				return
	// 			}

	// 			time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	// 			retries++

	// 		}

	// 		// for metricName, metricValue := range metrics {
	// 		// 	var metricType string
	// 		// 	if metricName == constants.PollCount {
	// 		// 		metricType = constants.Counter
	// 		// 	} else {
	// 		// 		metricType = constants.Gauge
	// 		// 	}
	// 		// 	err := sender.SendMetric(serverAddress, metricType, metricName, metricValue)

	// 		// 	if err != nil {
	// 		// 		fmt.Printf("Ошибка при отправке метрики %s: %s\n", metricName, err)
	// 		// 	}
	// 		// }
	// 	}
	// }
}
