package storage

import (
	"fmt"
	"strconv"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
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

func (ms *MemStorage) SetGauge(key string, value float64) {
	ms.gauge[key] = value
}

func (ms *MemStorage) StrValueToFloat(str string) (value float64, err error) {
	value, err = strconv.ParseFloat(str, 64)
	return
}

func (ms *MemStorage) StrValueToInt(str string) (value int64, err error) {
	value, err = strconv.ParseInt(str, 10, 64)
	return
}

func (ms *MemStorage) SetCounter(key string, value int64) {
	_, ok := ms.counter[key]
	if ok {
		ms.counter[key] += value

	} else {
		ms.counter[key] = value
	}
}

func (ms *MemStorage) GetGauge(key string) (value float64, err error) {
	value, ok := ms.gauge[key]
	if ok != true {
		err = fmt.Errorf("Значение метрики %s типа gauge не найдено", key)
	}
	return
}

func (ms *MemStorage) GetCounter(key string) (value int64, err error) {
	value, ok := ms.counter[key]
	if ok != true {
		err = fmt.Errorf("Значение метрики %s типа Counter не найдено", key)
	}
	return
}

func (ms *MemStorage) GetAllGauge() map[string]float64 {
	return ms.gauge[]
} 

func (ms *MemStorage) GetAllCounter() map[string]int64 {
	return ms.counter[]
} 

func (ms *MemStorage) GetAll() result map[string]interface {
	result = ms.gauge[]
	for key, value := range ms.counter {
		result[key] = value
	}
	return
}
