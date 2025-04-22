package agent

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var pollCount int64

func (agent *Agent) collectRuntimeMetrics() {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	agent.mutex.Lock()
	agent.metrics["Alloc"] = memStats.Alloc
	agent.metrics["BuckHashSys"] = memStats.BuckHashSys
	agent.metrics["Frees"] = memStats.Frees
	agent.metrics["GCCPUFraction"] = memStats.GCCPUFraction
	agent.metrics["GCSys"] = memStats.GCSys
	agent.metrics["HeapAlloc"] = memStats.HeapAlloc
	agent.metrics["HeapIdle"] = memStats.HeapIdle
	agent.metrics["HeapInuse"] = memStats.HeapInuse
	agent.metrics["HeapObjects"] = memStats.HeapObjects
	agent.metrics["HeapReleased"] = memStats.HeapReleased
	agent.metrics["HeapSys"] = memStats.HeapSys
	agent.metrics["LastGC"] = memStats.LastGC
	agent.metrics["Lookups"] = memStats.Lookups
	agent.metrics["MCacheInuse"] = memStats.MCacheInuse
	agent.metrics["MCacheSys"] = memStats.MCacheSys
	agent.metrics["MSpanInuse"] = memStats.MSpanInuse
	agent.metrics["MSpanSys"] = memStats.MSpanSys
	agent.metrics["Mallocs"] = memStats.Mallocs
	agent.metrics["NextGC"] = memStats.NextGC
	agent.metrics["NumForcedGC"] = memStats.NumForcedGC
	agent.metrics["NumGC"] = memStats.NumGC
	agent.metrics["OtherSys"] = memStats.OtherSys
	agent.metrics["PauseTotalNs"] = memStats.PauseTotalNs
	agent.metrics["StackInuse"] = memStats.StackInuse
	agent.metrics["StackSys"] = memStats.StackSys
	agent.metrics["Sys"] = memStats.Sys
	agent.metrics["TotalAlloc"] = memStats.TotalAlloc

	pollCount++
	agent.metrics["PollCount"] = pollCount

	agent.metrics["RandomValue"] = rand.Float64()

	agent.mutex.Unlock()

	slog.Info(fmt.Sprintf("runtime metrics saved: PollCount = %v", pollCount))
}

func (agent *Agent) setPollCountInitial() {
	pollCount = 0
}

func (agent *Agent) ResetMetricsStorage() {
	agent.metrics = make(map[string]interface{})
	slog.Info("metrics map reseted")
}

func (agent *Agent) collectAdditionalMetrics() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return
	}

	agent.mutex.Lock()

	agent.metrics["TotalMemory"] = vm.Total
	agent.metrics["FreeMemory"] = vm.Free

	for i, percent := range cpuPercent {
		agent.metrics[fmt.Sprintf("CPUutilization%d", i+1)] = percent
	}

	agent.mutex.Unlock()

	slog.Info(fmt.Sprintf("additional metrics saved: PollCount = %v", pollCount))

}
