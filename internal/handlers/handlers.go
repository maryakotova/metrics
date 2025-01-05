package handlers

import (
	"net/http"
	"strings"

	. "github.com/maryakotova/metrics/internal/storage"
)

type Server struct {
	metrics MemStorage
}

func NewServer(metrics *MemStorage) *Server {
	return &Server{metrics: *metrics}
}

func (server *Server) HandleMetricUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		// http.Error(res, "Невозможно обновить метрику(недостаточно параметров)", http.StatusMethodNotAllowed)
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/plain")

	parcedURL := strings.Split(req.URL.Path, "/")
	if len(parcedURL) < 5 {
		http.Error(res, "Невозможно обновить метрику(недостаточно параметров)", http.StatusNotFound)
		return
	} else if parcedURL[1] != "update" {
		http.Error(res, "Невозможно обновить метрику(где update?) "+parcedURL[2], http.StatusNotFound)
		return
	} else if len(parcedURL) == 6 && parcedURL[6] != "" || len(parcedURL) > 6 {
		http.Error(res, "Невозможно обновить метрику(слишком много параметров)", http.StatusNotFound)
	}

	metricType := parcedURL[2]
	metricName := parcedURL[3]
	metricValue := parcedURL[4]

	if metricName == "" || metricValue == "" {
		http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
		return
	}

	switch metricType {
	case "gauge":
		gaugeValue, err := server.metrics.StrValueToFloat(metricValue)
		if err != nil {
			http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
			return
		}
		server.metrics.SetGauge(metricName, gaugeValue)

	case "counter":
		counterValue, err := server.metrics.StrValueToInt(metricValue)
		if err != nil {
			http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
			return
		}
		server.metrics.SetCounter(metricName, counterValue)

	default:
		http.Error(res, "Неверный формат для обновления метрик (неверное имя)", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)

}
