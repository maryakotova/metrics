package storage

// type MetricType string

import (
	"fmt"
	"strconv"
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

func (ms *MemStorage) SetGauge(key string, value string) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}
	floatValue, err := ms.strValueToFloat(value)
	if err != nil {
		return err
	}

	ms.gauge[key] = floatValue
	return
}

func (ms *MemStorage) strValueToInt(str string) (value int64, err error) {
	value, err = strconv.ParseInt(str, 10, 64)
	return
}

func (ms *MemStorage) SetCounter(key string, value string) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}
	intValue, err := ms.strValueToInt(value)
	if err != nil {
		return err
	}

	_, ok := ms.counter[key]
	if ok {
		ms.counter[key] += intValue

	} else {
		ms.counter[key] = intValue
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
