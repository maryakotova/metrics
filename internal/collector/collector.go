package collector

import (
	"math/rand/v2"
	"runtime"
)

var pollCount int64

// func CollectMetrics() []models.MetricsForSend {
// 	metrics := []models.MetricsForSend{}

// 	memStats := new(runtime.MemStats)
// 	runtime.ReadMemStats(memStats)

// 	var value float64

// 	value = float64(memStats.Alloc)
// 	metrics = append(metrics, models.MetricsForSend{ID: "Alloc", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.BuckHashSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "BuckHashSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.Frees)
// 	metrics = append(metrics, models.MetricsForSend{ID: "Frees", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.GCCPUFraction)
// 	metrics = append(metrics, models.MetricsForSend{ID: "GCCPUFraction", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.GCSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "GCSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.HeapAlloc)
// 	metrics = append(metrics, models.MetricsForSend{ID: "HeapAlloc", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.HeapIdle)
// 	metrics = append(metrics, models.MetricsForSend{ID: "HeapIdle", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.HeapInuse)
// 	metrics = append(metrics, models.MetricsForSend{ID: "HeapInuse", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.HeapObjects)
// 	metrics = append(metrics, models.MetricsForSend{ID: "HeapObjects", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.HeapReleased)
// 	metrics = append(metrics, models.MetricsForSend{ID: "HeapReleased", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.HeapSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "HeapSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.LastGC)
// 	metrics = append(metrics, models.MetricsForSend{ID: "LastGC", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.Lookups)
// 	metrics = append(metrics, models.MetricsForSend{ID: "Lookups", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.MCacheInuse)
// 	metrics = append(metrics, models.MetricsForSend{ID: "MCacheInuse", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.MCacheSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "MCacheSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.MSpanInuse)
// 	metrics = append(metrics, models.MetricsForSend{ID: "MSpanInuse", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.MSpanSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "MSpanSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.Mallocs)
// 	metrics = append(metrics, models.MetricsForSend{ID: "Mallocs", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.NextGC)
// 	metrics = append(metrics, models.MetricsForSend{ID: "NextGC", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.NumForcedGC)
// 	metrics = append(metrics, models.MetricsForSend{ID: "NumForcedGC", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.NumGC)
// 	metrics = append(metrics, models.MetricsForSend{ID: "NumGC", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.OtherSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "OtherSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.PauseTotalNs)
// 	metrics = append(metrics, models.MetricsForSend{ID: "PauseTotalNs", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.StackInuse)
// 	metrics = append(metrics, models.MetricsForSend{ID: "StackInuse", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.StackSys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "StackSys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.Sys)
// 	metrics = append(metrics, models.MetricsForSend{ID: "Sys", MType: constants.Gauge, Value: value})

// 	value = float64(memStats.TotalAlloc)
// 	metrics = append(metrics, models.MetricsForSend{ID: "TotalAlloc", MType: constants.Gauge, Value: value})

// 	value = rand.Float64()
// 	metrics = append(metrics, models.MetricsForSend{ID: "RandomValue", MType: constants.Gauge, Value: value})

// 	pollCount++
// 	metrics = append(metrics, models.MetricsForSend{ID: "PollCount", MType: constants.Counter, Delta: pollCount})

// 	return metrics
// }

func CollectMetrics() map[string]interface{} {
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

func SetPollCountInitial() {
	pollCount = 0
}
