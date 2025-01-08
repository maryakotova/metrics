package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var (
	serverAddress  string
	pollInterval   int64
	reportInterval int64
)

var pollCount int64

func collectMetrics() map[string]interface{} {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	metrics := make(map[string]interface{})

	metrics["Alloc"] = memStats.Alloc
	metrics["BuckHashSys"] = memStats.BuckHashSys
	metrics["Frees"] = memStats.Frees
	metrics["GCCPUFraction"] = memStats.GCCPUFraction
	metrics["GCSys"] = memStats.GCSys
	metrics["HeapAlloc"] = memStats.HeapAlloc
	metrics["HeapIdle"] = memStats.HeapIdle
	metrics["HeapInuse"] = memStats.HeapInuse
	metrics["HeapObjects"] = memStats.HeapObjects
	metrics["HeapReleased"] = memStats.HeapReleased
	metrics["HeapSys"] = memStats.HeapSys
	metrics["LastGC"] = memStats.LastGC
	metrics["Lookups"] = memStats.Lookups
	metrics["MCacheInuse"] = memStats.MCacheInuse
	metrics["MCacheSys"] = memStats.MCacheSys
	metrics["MSpanInuse"] = memStats.MSpanInuse
	metrics["MSpanSys"] = memStats.MSpanSys
	metrics["Mallocs"] = memStats.Mallocs
	metrics["NextGC"] = memStats.NextGC
	metrics["NumForcedGC"] = memStats.NumForcedGC
	metrics["NumGC"] = memStats.NumGC
	metrics["OtherSys"] = memStats.OtherSys
	metrics["PauseTotalNs"] = memStats.PauseTotalNs
	metrics["StackInuse"] = memStats.StackInuse
	metrics["StackSys"] = memStats.StackSys
	metrics["Sys"] = memStats.Sys
	metrics["TotalAlloc"] = memStats.TotalAlloc

	pollCount++
	metrics["PollCount"] = pollCount

	metrics["RandomValue"] = rand.Float64()

	return metrics
}

func sendMetric(metricType string, metricName string, metricValue interface{}) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", serverAddress, metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending metric:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Sent metric: %s/%s/%v\n", metricType, metricName, metricValue)
	return err
}

func main() {

	parseFlags()

	n := int64(0)

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		n += pollInterval

		metrics := collectMetrics()

		if reportInterval == n {
			n = 0
			pollCount = 0

			for metricName, metricValue := range metrics {
				var metricType string
				if metricName == "PollCount" {
					metricType = "counter"
				} else {
					metricType = "gauge"
				}
				err := sendMetric(metricType, metricName, metricValue)
				if err != nil {
					fmt.Printf("Ошибка при отправке метрики %s: %s\n", metricName, err)
				}
			}
		}
	}
}
