package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"github.com/maryakotova/metrics/internal/models"
)

const tplPath string = "templates/metrics.html"

type DataStorage interface {
	SetGauge(key string, value float64) (err error)
	SetCounter(key string, value int64) (err error)
	GetAll() map[string]interface{}
	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
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
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "Content-Type: application/json")

	var request models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {

		metricType := req.PathValue("metricType")
		metricName := req.PathValue("metricName")
		metricValue := req.PathValue("metricValue")

		if metricName == "" {
			http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
			return
		}

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
				return
			}
			err = server.metrics.SetGauge(metricName, value)
			if err != nil {
				http.Error(res, "Неверный формат значения для обновления метрики Gauge", http.StatusBadRequest)
				return
			}

		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
				return
			}
			err = server.metrics.SetCounter(metricName, value)
			if err != nil {
				http.Error(res, "Неверный формат значения для обновления метрик Counter", http.StatusBadRequest)
				return
			}

		default:
			http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
			return
		}

	} else {

		if request.ID == "" {
			http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
			return
		}

		switch request.MType {
		case "gauge":
			if err := server.metrics.SetGauge(request.ID, *request.Value); err != nil {
				http.Error(res, "Ошибка при обновлении метрики Gauge", http.StatusBadRequest)
			}
		case "counter":
			if err := server.metrics.SetCounter(request.ID, *request.Delta); err != nil {
				http.Error(res, "Ошибка при обновлении метрики Counter", http.StatusBadRequest)
			}
		default:
			http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
			return
		}
	}

	res.WriteHeader(http.StatusOK)

}

func (server *Server) HandleGetOneMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	responce := models.Metrics{
		ID:    request.ID,
		MType: request.MType,
	}

	switch request.MType {
	case "gauge":
		gaugeValue, err := server.metrics.GetGauge(request.ID)
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		responce.Value = &gaugeValue
	case "counter":
		counterValue, err := server.metrics.GetCounter(request.ID)
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		responce.Delta = &counterValue
	default:
		http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(res)
	if err := enc.Encode(responce); err != nil {
		http.Error(res, "Ошибка при заполнении ответа", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)

}

func (server *Server) HandleGetAllMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/html")

	data := struct {
		IntMap   map[string]int64
		FloatMap map[string]float64
	}{
		IntMap:   server.metrics.GetAllCounter(),
		FloatMap: server.metrics.GetAllGauge(),
	}

	tmpl, err := template.ParseFiles(tplPath)
	if err != nil {
		http.Error(res, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(res, data)
	if err != nil {
		http.Error(res, "Error executing template", http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)
}
