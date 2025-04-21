package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"metrics/internal/authsign"
	"metrics/internal/config"
	"metrics/internal/constants"
	"metrics/internal/controller"
	"metrics/internal/filetransfer"
	"metrics/internal/models"

	"go.uber.org/zap"
)

// const tplPath string = "./templates/metrics.html"
const tplPath string = "templates/metrics.html"

// type DataStorage interface {
// 	SetGauge(ctx context.Context, key string, value float64) (err error)
// 	SetCounter(ctx context.Context, key string, value *int64) (err error)
// 	SaveMetrics(ctx context.Context, metrics []models.Metrics) (err error)
// 	GetAllGauge(ctx context.Context) map[string]float64
// 	GetAllCounter(ctx context.Context) map[string]int64
// 	GetGauge(ctx context.Context, key string) (value float64, err error)
// 	GetCounter(ctx context.Context, key string) (value int64, err error)
// 	CheckConnection(ctx context.Context) (err error)
// }

// ----------------------------------------------------------------------
//fileWriter должен остаться в сервере? или перейти в контроллер?
// ----------------------------------------------------------------------

type Server struct {
	config     *config.Config
	fileWriter *filetransfer.FileWriter
	logger     *zap.Logger
	controller *controller.Controller
}

func NewServer(cfg *config.Config, fileWriter *filetransfer.FileWriter, logger *zap.Logger, controller *controller.Controller) *Server {
	return &Server{
		config:     cfg,
		fileWriter: fileWriter,
		logger:     logger,
		controller: controller,
	}
}

func (server *Server) HandleMetricUpdateViaJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	receivedHash := req.Header.Get(constants.HeaderSig)
	if receivedHash != "" {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := buf.Bytes()
		if !authsign.VerifySig(receivedHash, []byte(body), []byte(server.config.SecretKey)) {
			err = fmt.Errorf("invalid hash")
			server.logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// ----------------------------------------------------------------------------------------------------------------------------------------------
	// При наличии ключа на этапе формирования ответа сервер должен вычислять хеш и передавать его в HTTP-заголовке ответа с именем HashSHA256.
	// Сервер должен отправлять тот же самый хэш? или новый в зависимости от ответа?
	// ----------------------------------------------------------------------------------------------------------------------------------------------

	var request models.Metrics

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		err = fmt.Errorf("ошибка в JSON: %w", err)
		server.logger.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	statusCode, err := server.controller.UpdateMetric(req.Context(), request)
	if err != nil {
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		http.Error(res, err.Error(), statusCode)
		return
	}

	responce := request

	res.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(res)
	if err := enc.Encode(responce); err != nil {
		err = fmt.Errorf("ошибка при заполнении ответа: %w", err)
		server.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

	if server.config.IsSyncStore() {
		server.fileWriter.WriteMetrics(responce)
	}
}

func (server *Server) HandleMetricUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	receivedHash := req.Header.Get(constants.HeaderSig)
	if receivedHash != "" {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := buf.Bytes()
		if !authsign.VerifySig(receivedHash, []byte(body), []byte(server.config.SecretKey)) {
			err = fmt.Errorf("invalid hash")
			server.logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")

	statusCode, err := server.controller.UpdateMetricFromString(req.Context(), metricType, metricName, &metricValue)
	if err != nil {
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		http.Error(res, err.Error(), statusCode)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	// ----------------------------------------------------------------------
	//как правильно заполнить responce в данной ситуации (value и delta) ?
	// ----------------------------------------------------------------------
	responce := models.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case constants.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err == nil {
			responce.Value = &value
		}
	case constants.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err == nil {
			responce.Delta = &value
		}
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(responce); err != nil {
		err = fmt.Errorf("ошибка при заполнении ответа: %w", err)
		server.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

	if server.config.IsSyncStore() {
		server.fileWriter.WriteMetrics(responce)
	}
}

func (server *Server) HandleGetOneMetricViaJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	receivedHash := req.Header.Get(constants.HeaderSig)
	if receivedHash != "" {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := buf.Bytes()
		if !authsign.VerifySig(receivedHash, []byte(body), []byte(server.config.SecretKey)) {
			err = fmt.Errorf("invalid hash")
			server.logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var request models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		err = fmt.Errorf("ошибка в JSON: %w", err)
		server.logger.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	responce := request

	statusCode, err := server.controller.GetOneMetric(req.Context(), &responce)
	if err != nil {
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		http.Error(res, err.Error(), statusCode)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(res)
	if err := enc.Encode(responce); err != nil {
		err = fmt.Errorf("ошибка при заполнении ответа: %w", err)
		server.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (server *Server) HandleGetOneMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/plain")

	data := models.Metrics{}
	data.ID = req.PathValue("metricName")
	data.MType = req.PathValue("metricType")

	statusCode, err := server.controller.GetOneMetric(req.Context(), &data)
	if err != nil {
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		http.Error(res, err.Error(), statusCode)
		return
	}

	switch data.MType {
	case constants.Gauge:
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatFloat(*data.Value, 'f', -1, 64)))
	case constants.Counter:
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatInt(*data.Delta, 10)))
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
	}{IntMap: server.controller.GetAllCounter(req.Context()), FloatMap: server.controller.GetAllGauge(req.Context())}

	tmpl, err := template.ParseFiles(tplPath)
	if err != nil {
		err = fmt.Errorf("error parsing template: %w", err)
		server.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(res, data)
	if err != nil {
		err = fmt.Errorf("error executing template: %w", err)
		server.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (server *Server) HandlePing(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Content-Type", "text/plain")

	if err := server.controller.CheckConnection(req.Context()); err != nil {
		server.logger.Error(err.Error())
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

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := buf.Bytes()

	receivedHash := req.Header.Get(constants.HeaderSig)
	if receivedHash != "" {
		if !authsign.VerifySig(receivedHash, []byte(body), []byte(server.config.SecretKey)) {
			err = fmt.Errorf("invalid hash")
			server.logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var request []models.Metrics

	// decoder := json.NewDecoder(req.Body)
	// if err := decoder.Decode(&request); err != nil {
	if err = json.Unmarshal([]byte(body), &request); err != nil {
		err = fmt.Errorf("ошибка в JSON: %w", err)
		server.logger.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(request) == 0 {
		http.Error(res, "Empty batch", http.StatusBadRequest)
		return
	}

	statusCode, err := server.controller.SaveMetrics(req.Context(), request)
	if err != nil {
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		err = fmt.Errorf("ошибка при сохранении: %w", err)
		server.logger.Error(err.Error())
		http.Error(res, err.Error(), statusCode)
		return
	}

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("metrics have been updated"))

}
