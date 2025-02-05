package app

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/maryakotova/metrics/internal/models"
// )

// func MetricUpdate(metricType string, metricName string, metricValue string) (responce models.Metrics, err error, status int) {

// 	if metricName == "" {
// 		err = fmt.Errorf("Невозможно обновить метрику(пустое имя или значение метрики)")
// 		status = http.StatusNotFound
// 		return
// 	}

// 	responce.ID = metricName
// 	responce.MType = metricType

// 	switch metricType {
// 	case "gauge":
// 		value, er := strconv.ParseFloat(metricValue, 64)
// 		if er != nil {
// 			err = fmt.Errorf("Неверный формат значения для обновления метрики Gauge")
// 			status = http.StatusBadRequest
// 			return
// 		}
// 		er = server.metrics.SetGauge(metricName, value)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
// 			return
// 		}
// 		responce.Value = &value

// 	case "counter":
// 		value, err := strconv.ParseInt(metricValue, 10, 64)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
// 			return
// 		}
// 		err = server.metrics.SetCounter(metricName, value)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
// 			return
// 		}
// 		responce.Delta = &value

// 	default:
// 		http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
// 		return
// 	}

// }
