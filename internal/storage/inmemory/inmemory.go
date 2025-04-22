package inmemory

import (
	"context"
	"fmt"
	"sync"

	"metrics/internal/constants"
	"metrics/internal/filetransfer"
	"metrics/internal/models"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	m       sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) SetGauge(ctx context.Context, key string, value float64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return err
	}
	ms.m.Lock()
	ms.gauge[key] = value
	ms.m.Unlock()
	return
}

func (ms *MemStorage) SetCounter(ctx context.Context, key string, value *int64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return err
	}

	ms.m.Lock()
	val, ok := ms.counter[key]
	if ok {
		ms.counter[key] += *value
		*value += val

	} else {
		ms.counter[key] = *value
	}
	ms.m.Unlock()
	return
}

func (ms *MemStorage) SetCounterFromFile(ctx context.Context, key string, value int64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return err
	}
	ms.m.Lock()
	ms.counter[key] = value
	ms.m.Unlock()
	return
}

func (ms *MemStorage) GetGauge(ctx context.Context, key string) (value float64, err error) {
	ms.m.Lock()
	value, ok := ms.gauge[key]
	ms.m.Unlock()
	if !ok {
		err = fmt.Errorf("значение метрики %s типа gauge не найдено", key)
		return
	}
	return
}

func (ms *MemStorage) GetCounter(ctx context.Context, key string) (value int64, err error) {
	ms.m.Lock()
	value, ok := ms.counter[key]
	ms.m.Unlock()
	if !ok {
		err = fmt.Errorf("значение метрики %s типа Counter не найдено", key)
	}
	return
}

func (ms *MemStorage) GetAllGauge(ctx context.Context) map[string]float64 {
	return ms.gauge

}

func (ms *MemStorage) GetAllCounter(ctx context.Context) map[string]int64 {
	return ms.counter
}

func (ms *MemStorage) GetAll(ctx context.Context) map[string]interface{} {

	allMetrics := make(map[string]interface{})

	ms.m.Lock()
	for key, value := range ms.gauge {
		allMetrics[key] = value
	}

	for key, value := range ms.counter {
		allMetrics[key] = value
	}
	ms.m.Unlock()

	return allMetrics
}

func (ms *MemStorage) GetAllMetricsInJSON() []models.Metrics {

	metrics := []models.Metrics{}

	ms.m.Lock()
	for key, value := range ms.gauge {

		metric := models.Metrics{
			ID:    key,
			MType: constants.Gauge,
			Value: &value}
		metrics = append(metrics, metric)
	}

	for key, value := range ms.counter {
		metric := models.Metrics{
			ID:    key,
			MType: constants.Counter,
			Delta: &value}
		metrics = append(metrics, metric)
	}
	ms.m.Unlock()

	return metrics
}

func (ms *MemStorage) UploadData(filePath string) {

	fileReader, err := filetransfer.NewFileReader(filePath)
	if err != nil {
		panic(err)
	}

	metrics, err := fileReader.ReadMetrics()
	if err != nil {
		return
	}

	defer fileReader.Close()

	if len(metrics) > 0 {
		for _, metric := range metrics {
			switch metric.MType {
			case constants.Gauge:
				ms.SetGauge(context.Background(), metric.ID, *metric.Value)
			case constants.Counter:
				ms.SetCounterFromFile(context.Background(), metric.ID, *metric.Delta)
			}
		}
	}
}

func (ms *MemStorage) CheckConnection(ctx context.Context) (err error) {
	return fmt.Errorf("in-memory storage is used")
}

func (ms *MemStorage) SaveMetrics(ctx context.Context, metrics []models.Metrics) (err error) {
	for _, metric := range metrics {
		var zero float64 = 0
		switch metric.MType {
		case constants.Gauge:
			if metric.Value == nil {
				metric.Value = &zero
			}
			if err = ms.SetGauge(ctx, metric.ID, *metric.Value); err != nil {
				return err
			}
		case constants.Counter:
			if err = ms.SetCounter(ctx, metric.ID, metric.Delta); err != nil {
				return err
			}
		default:
			return fmt.Errorf("неверный формат для обновления метрик (недопустимый тип): %s", metric.MType)
		}
	}
	return
}
