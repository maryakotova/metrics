package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/filetransfer"
	"github.com/maryakotova/metrics/internal/models"
)

const tplPath string = "./templates/metrics.html"

type DataStorage interface {
	SetGauge(ctx context.Context, key string, value float64) (err error)
	SetCounter(ctx context.Context, key string, value int64) (err error)
	SaveMetrics(ctx context.Context, metrics []models.Metrics) (err error)
	//GetAll(ctx context.Context) map[string]interface{}
	GetAllGauge(ctx context.Context) map[string]float64
	GetAllCounter(ctx context.Context) map[string]int64
	GetGauge(ctx context.Context, key string) (value float64, err error)
	GetCounter(ctx context.Context, key string) (value int64, err error)
	CheckConnection(ctx context.Context) (err error)
}

type Server struct {
	storage       DataStorage
	syncFileWrite bool
	fileWriter    *filetransfer.FileWriter
}

func NewServer(storage DataStorage, syncFileWrite bool, fileWriter *filetransfer.FileWriter) *Server {
	return &Server{
		storage:       storage,
		syncFileWrite: syncFileWrite,
		fileWriter:    fileWriter,
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
		if err := server.storage.SetGauge(req.Context(), request.ID, *request.Value); err != nil {
			http.Error(res, "Ошибка при обновлении метрики Gauge", http.StatusBadRequest)
			return
		}
		responce.Value = request.Value
	case constants.Counter:
		if err := server.storage.SetCounter(req.Context(), request.ID, *request.Delta); err != nil {
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
		err = server.storage.SetGauge(req.Context(), metricName, value)
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
		err = server.storage.SetCounter(req.Context(), metricName, value)
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
		gaugeValue, err := server.storage.GetGauge(req.Context(), request.ID)
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		responce.Value = &gaugeValue
	case constants.Counter:
		counterValue, err := server.storage.GetCounter(req.Context(), request.ID)
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
		metricValue, err := server.storage.GetGauge(req.Context(), req.PathValue("metricName"))
		if err != nil {
			http.Error(res, "Запрос неизвестной метрики", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatFloat(metricValue, 'f', -1, 64)))

	case constants.Counter:
		metricValue, err := server.storage.GetCounter(req.Context(), req.PathValue("metricName"))
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
	}{IntMap: server.storage.GetAllCounter(req.Context()), FloatMap: server.storage.GetAllGauge(req.Context())}

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

	if err := server.storage.CheckConnection(req.Context()); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Connection is successful"))
}

func (server *Server) HandleMetricUpdates(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request []models.Metrics

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(request) == 0 {
		http.Error(res, "Empty batch", http.StatusBadRequest)
		return
	}

	if err := server.storage.SaveMetrics(req.Context(), request); err != nil {
		http.Error(res, "error when saving to DB", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("metrics have been updated"))

}
