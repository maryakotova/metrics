package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/filetransfer"
	"github.com/maryakotova/metrics/internal/models"
)

const tplPath string = "./templates/metrics.html"

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
	metrics       DataStorage
	syncFileWrite bool
	fileWriter    *filetransfer.FileWriter
	db            *sql.DB
}

func NewServer(metrics DataStorage, syncFileWrite bool, fileWriter *filetransfer.FileWriter, db *sql.DB) *Server {
	return &Server{
		metrics:       metrics,
		syncFileWrite: syncFileWrite,
		fileWriter:    fileWriter,
		db:            db,
	}
}

func (server *Server) HandleMetricUpdateViaJSON(res http.ResponseWriter, req *http.Request) {
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

	if request.ID == "" {
		http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
		return
	}

	responce := models.Metrics{
		ID:    request.ID,
		MType: request.MType,
	}

	switch request.MType {
	case constants.Gauge:
		if err := server.metrics.SetGauge(request.ID, *request.Value); err != nil {
			http.Error(res, "Ошибка при обновлении метрики Gauge", http.StatusBadRequest)
			return
		}
		responce.Value = request.Value
	case constants.Counter:
		if err := server.metrics.SetCounter(request.ID, *request.Delta); err != nil {
			http.Error(res, "Ошибка при обновлении метрики Counter", http.StatusBadRequest)
			return
		}
		responce.Delta = request.Delta
	default:
		http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "Content-Type: application/json")

	enc := json.NewEncoder(res)
	if err := enc.Encode(responce); err != nil {
		http.Error(res, "Ошибка при заполнении ответа", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	// fmt.Printf("Responce: %v\n", responce)

	if server.syncFileWrite {
		server.fileWriter.WriteMetric(&responce)
	}
}

func (server *Server) HandleMetricUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "Content-Type: application/json")

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")

	if metricName == "" {
		http.Error(res, "Невозможно обновить метрику(пустое имя или значение метрики)", http.StatusNotFound)
		return
	}

	responce := models.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case constants.Gauge:
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
		responce.Value = &value

	case constants.Counter:
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
		responce.Delta = &value

	default:
		http.Error(res, "Неверный формат для обновления метрик (неверный тип)", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(responce); err != nil {
		http.Error(res, "Ошибка при заполнении ответа", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	if server.syncFileWrite {
		server.fileWriter.WriteMetric(&responce)
	}
}

func (server *Server) HandleGetOneMetricViaJSON(res http.ResponseWriter, req *http.Request) {
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
	case constants.Gauge:
		gaugeValue, err := server.metrics.GetGauge(request.ID)
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		responce.Value = &gaugeValue
	case constants.Counter:
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
	fmt.Printf("Responce: %v\n", responce)

}

func (server *Server) HandleGetOneMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/plain")

	// switch chi.URLParam(req, "type") { !!!Почему не работает??????
	switch req.PathValue("metricType") {
	case constants.Gauge:
		metricValue, err := server.metrics.GetGauge(req.PathValue("metricName"))
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatFloat(metricValue, 'f', -1, 64)))

	case constants.Counter:
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

func (server *Server) HandlePing(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/plain")

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if err := server.db.PingContext(ctx); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

}
