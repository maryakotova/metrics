package storage

import (
	"fmt"
	"strconv"

	"github.com/maryakotova/metrics/internal/models"
)

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) strValueToFloat(str string) (value float64, err error) {
	value, err = strconv.ParseFloat(str, 64)
	return
}

func (ms *MemStorage) SetGauge(key string, value float64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}

	ms.gauge[key] = value
	return
}

func (ms *MemStorage) strValueToInt(str string) (value int64, err error) {
	value, err = strconv.ParseInt(str, 10, 64)
	return
}

func (ms *MemStorage) SetCounter(key string, value int64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}
	_, ok := ms.counter[key]
	if ok {
		ms.counter[key] += value

	} else {
		ms.counter[key] = value
	}
	return
}

func (ms *MemStorage) GetGauge(key string) (value float64, err error) {

	value, ok := ms.gauge[key]
	if !ok {
		err = fmt.Errorf("значение метрики %s типа gauge не найдено", key)
	}
	return
}

func (ms *MemStorage) GetCounter(key string) (value int64, err error) {
	value, ok := ms.counter[key]
	if !ok {
		err = fmt.Errorf("значение метрики %s типа Counter не найдено", key)
	}
	return
}

func (ms *MemStorage) GetAllGauge() map[string]float64 {
	return ms.gauge
}

func (ms *MemStorage) GetAllCounter() map[string]int64 {
	return ms.counter
}

func (ms *MemStorage) GetAll() map[string]interface{} {

	allMetrics := make(map[string]interface{})

	for key, value := range ms.gauge {
		allMetrics[key] = value
	}

	for key, value := range ms.counter {
		allMetrics[key] = value
	}

	return allMetrics
}

func (ms *MemStorage) GetAllMetricsInJSON() []models.Metrics {

	metrics := []models.Metrics{}

	for key, value := range ms.gauge {

		metric := models.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value}
		metrics = append(metrics, metric)
	}

	for key, value := range ms.counter {
		metric := models.Metrics{
			ID:    key,
			MType: "counter",
			Delta: &value}
		metrics = append(metrics, metric)
	}

	return metrics
}
