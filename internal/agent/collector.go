package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// var pollCount int64

func (agent *Agent) collectRuntimeMetrics() map[string]interface{} {

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

	// agent.setPollCountInitial()

	// pollCount++
	metrics["PollCount"] = 1

	metrics["RandomValue"] = rand.Float64()

	return metrics
}

// func (agent *Agent) setPollCountInitial() {
// 	pollCount = 0
// }

func (agent *Agent) collectAdditionalMetrics() (map[string]interface{}, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]interface{})

	metrics["TotalMemory"] = vm.Total
	metrics["FreeMemory"] = vm.Free
	metrics["CPUutilization1"] = cpuPercent

	return metrics, nil
}
