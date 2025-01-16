package logic

// import "net/http"

// func UpdateMetric(metricType string, metricName string, metricValue string) (errorText string, code int) {
// 	if metricName == "" {
// 		errorText = "Невозможно обновить метрику(пустое имя или значение метрики"
// 		code = http.StatusNotFound
// 		return
// 	}

// 	switch metricType {
// 	case "gauge":
// 		err := server.metrics.SetGauge(metricName, metricValue)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
// 			return
// 		}

// 	case "counter":
// 		err := server.metrics.SetCounter(metricName, metricValue)
// 		if err != nil {
// 			http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
// 			return
// 		}

// 	default:
// 		http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
// 		return
// 	}

// 	return
// }
