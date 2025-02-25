package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var pollCount int64

// func (agent *Agent) collectRuntimeMetrics() {
// 	memStats := new(runtime.MemStats)
// 	runtime.ReadMemStats(memStats)

// 	agent.metrics["Alloc"] = memStats.Alloc
// 	agent.metrics["BuckHashSys"] = memStats.BuckHashSys
// 	agent.metrics["Frees"] = memStats.Frees
// 	agent.metrics["GCCPUFraction"] = memStats.GCCPUFraction
// 	agent.metrics["GCSys"] = memStats.GCSys
// 	agent.metrics["HeapAlloc"] = memStats.HeapAlloc
// 	agent.metrics["HeapIdle"] = memStats.HeapIdle
// 	agent.metrics["HeapInuse"] = memStats.HeapInuse
// 	agent.metrics["HeapObjects"] = memStats.HeapObjects
// 	agent.metrics["HeapReleased"] = memStats.HeapReleased
// 	agent.metrics["HeapSys"] = memStats.HeapSys
// 	agent.metrics["LastGC"] = memStats.LastGC
// 	agent.metrics["Lookups"] = memStats.Lookups
// 	agent.metrics["MCacheInuse"] = memStats.MCacheInuse
// 	agent.metrics["MCacheSys"] = memStats.MCacheSys
// 	agent.metrics["MSpanInuse"] = memStats.MSpanInuse
// 	agent.metrics["MSpanSys"] = memStats.MSpanSys
// 	agent.metrics["Mallocs"] = memStats.Mallocs
// 	agent.metrics["NextGC"] = memStats.NextGC
// 	agent.metrics["NumForcedGC"] = memStats.NumForcedGC
// 	agent.metrics["NumGC"] = memStats.NumGC
// 	agent.metrics["OtherSys"] = memStats.OtherSys
// 	agent.metrics["PauseTotalNs"] = memStats.PauseTotalNs
// 	agent.metrics["StackInuse"] = memStats.StackInuse
// 	agent.metrics["StackSys"] = memStats.StackSys
// 	agent.metrics["Sys"] = memStats.Sys
// 	agent.metrics["TotalAlloc"] = memStats.TotalAlloc

// 	pollCount++
// 	agent.metrics["PollCount"] = pollCount

// 	agent.metrics["RandomValue"] = rand.Float64()
// }

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

func (agent *Agent) setPollCountInitial() {
	pollCount = 0
}

// func (agent *Agent) collectAdditionalMetrics() {
// 	vm, err := mem.VirtualMemory()
// 	if err != nil {
// 		return
// 	}

// 	cpuPercent, err := cpu.Percent(0, true)
// 	if err != nil {
// 		return
// 	}

// 	agent.metrics["TotalMemory"] = vm.Total
// 	agent.metrics["FreeMemory"] = vm.Free
// 	agent.metrics["CPUutilization1"] = cpuPercent

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
