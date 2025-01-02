package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	pollInterval   int64 = 2
	reportInterval int64 = 10
)

var pollCont int64 = 0

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

	pollCont++
	metrics["PollCount"] = pollCont

	metrics["RandomValue"] = rand.Float64()

	return metrics

}

func sendMetric(metricType string, metricName string, metricValue interface{}) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	_, err = client.Do(req)

	// resp, err := client.Do(req)
	// // ===============
	// fmt.Printf("Status: %s\r\n", resp.Status)
	// fmt.Printf("Header ===============\r\n")
	// for k, v := range resp.Header {
	// 	fmt.Printf("%s: %v\r\n", k, v)
	// }
	// fmt.Printf("ContentLength: %v\n", resp.ContentLength)
	// fmt.Printf("Body: %v\n", resp.Body)
	// fmt.Printf("Query parameters ===============\r\n")

	// // ===============

	if err != nil {
		fmt.Println("Error sending metric:", err)
		return err
	}
	fmt.Printf("Sent metric: %s/%s/%v\n", metricType, metricName, metricValue)
	return err
}

func main() {

	n := int64(0)

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		n += pollInterval

		metrics := collectMetrics()

		if reportInterval == n {
			n = 0

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
				} else {
					// fmt.Printf("Метрика %s успешно отправлена \n", metricName)
				}
			}
		}
	}
}
