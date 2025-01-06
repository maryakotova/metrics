package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type DataStorage interface {
	SetGauge(key string, value string) (err error)
	SetCounter(key string, value string) (err error)
	GetAll() map[string]interface{}
	GetGauge(key string) (value float64, err error)
	GetCounter(key string) (value int64, err error)
}

type Server struct {
	metrics DataStorage
}

func NewServer(metrics DataStorage) *Server {
	return &Server{metrics: metrics}
}

func (server *Server) HandleMetricUpdate(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")

	if metricName == "" {
		http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
		return
	}

	switch metricType {
	case "gauge":
		err := server.metrics.SetGauge(metricName, metricValue)
		if err != nil {
			http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
			return
		}

	case "counter":
		err := server.metrics.SetCounter(metricName, metricValue)
		if err != nil {
			http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
			return
		}

	default:
		http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (server *Server) HandleGetOneMetric(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain")

	// switch chi.URLParam(req, "type") { !!!Почему не работает??????
	switch req.PathValue("metricType") {
	case "gauge":
		metricValue, err := server.metrics.GetGauge(req.PathValue("metricName"))
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatFloat(metricValue, 'f', -1, 64)))

	case "counter":
		metricValue, err := server.metrics.GetCounter(req.PathValue("metricName"))
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatInt(metricValue, 10)))

	case "":
		http.Error(res, "Тип метрики обязателен для заполнения", http.StatusBadRequest)
		return

	default:
		http.Error(res, "Указанный тип метрики не известен", http.StatusNotFound)
		return
	}
}

func (server *Server) HandleGetAllMetrics(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain")

	metricsList := server.metrics.GetAll()

	var sb strings.Builder
	for key, value := range metricsList {
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	res.Write([]byte(sb.String()))
}
