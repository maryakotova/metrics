package main

import (
	"net/http"

	. "github.com/maryakotova/metrics/internal/handlers"
	. "github.com/maryakotova/metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

// type MemStorage struct {
// 	gauge   map[string]float64
// 	counter map[string]int64
// }

// func NewMemStorage() *MemStorage {
// 	return &MemStorage{
// 		gauge:   make(map[string]float64),
// 		counter: make(map[string]int64),
// 	}
// }

// func (ms *MemStorage) SetGauge(key string, value float64) {
// 	ms.gauge[key] = value
// }

// func (ms *MemStorage) StrValueToFloat(str string) (value float64, err error) {
// 	value, err = strconv.ParseFloat(str, 64)
// 	return
// }

// func (ms *MemStorage) StrValueToInt(str string) (value int64, err error) {
// 	value, err = strconv.ParseInt(str, 10, 64)
// 	return
// }

// func (ms *MemStorage) SetCounter(key string, value int64) {
// 	_, ok := ms.counter[key]
// 	if ok {
// 		ms.counter[key] += value

// 	} else {
// 		ms.counter[key] = value
// 	}
// }

// func (ms *MemStorage) handleMetricUpdate(res http.ResponseWriter, req *http.Request) {
// 	if req.Method != http.MethodPost {
// 		// http.Error(res, "Невозможно обновить метрику(недостаточно параметров)", http.StatusMethodNotAllowed)
// 		res.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	res.Header().Set("Content-Type", "text/plain")

// 	parcedURL := strings.Split(req.URL.Path, "/")
// 	if len(parcedURL) < 5 {
// 		http.Error(res, "Невозможно обновить метрику(недостаточно параметров)", http.StatusNotFound)
// 		return
// 	} else if parcedURL[1] != "update" {
// 		http.Error(res, "Невозможно обновить метрику(где update?) "+parcedURL[2], http.StatusNotFound)
// 		return
// 	} else if len(parcedURL) == 6 && parcedURL[6] != "" || len(parcedURL) > 6 {
// 		http.Error(res, "Невозможно обновить метрику(слишком много параметров)", http.StatusNotFound)
// 	}

// 	metricType := parcedURL[2]
// 	metricName := parcedURL[3]
// 	metricValue := parcedURL[4]

// 	if metricName == "" || metricValue == "" {
// 		http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
// 		return
// 	}

// 	switch metricType {
// 	case "gauge":
// 		gaugeValue, err := ms.StrValueToFloat(metricValue)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
// 			return
// 		}
// 		ms.SetGauge(metricName, gaugeValue)

// 	case "counter":
// 		counterValue, err := ms.StrValueToInt(metricValue)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
// 			return
// 		}
// 		ms.SetCounter(metricName, counterValue)

// 	default:
// 		http.Error(res, "Неверный формат для обновления метрик (неверное имя)", http.StatusBadRequest)
// 		return
// 	}

// 	res.WriteHeader(http.StatusOK)

// }

// func handleBasic(res http.ResponseWriter, req *http.Request) {
// 	http.Error(res, "Неверная ссылка для обновления метрики", http.StatusNotFound)
// }

func main() {

	memStorage := NewMemStorage()
	server := NewServer(memStorage)

	router := chi.NewRouter()
	router.Use()

	router.Get("/", server.HandleGetAllMetrics)
	router.Get("/value/{metricType}/{metricName}", server.HandleGetOneMetric)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", server.HandleMetricUpdate)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
