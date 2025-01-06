package handlers

import (
	"net/http"
	"strconv"
	"text/template"
)

const tpl = `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Метрики</title>
		</head>
		<body>
			<h1>Метрики типа Counter</h1>
			<table border="1">
				<tr>
					<th>Key</th>
					<th>Value</th>
				</tr>
				{{range $key, $value := .IntMap}}
					<tr>
						<td>{{$key}}</td>
						<td>{{$value}}</td>
					</tr>
				{{end}}
			</table>

			<h1>Метрики типа Gauge</h1>
			<table border="1">
				<tr>
					<th>Key</th>
					<th>Value</th>
				</tr>
				{{range $key, $value := .FloatMap}}
					<tr>
						<td>{{$key}}</td>
						<td>{{$value}}</td>
					</tr>
				{{end}}
			</table>
		</body>
		</html>
		`

type DataStorage interface {
	SetGauge(key string, value string) (err error)
	SetCounter(key string, value string) (err error)
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

	res.Header().Set("Content-Type", "text/html")

	data := struct {
		IntMap   map[string]int64
		FloatMap map[string]float64
	}{
		IntMap:   server.metrics.GetAllCounter(),
		FloatMap: server.metrics.GetAllGauge(),
	}

	tmpl, err := template.New("webpage").Parse(tpl)
	if err != nil {
		http.Error(res, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(res, data)
	if err != nil {
		http.Error(res, "Error executing template", http.StatusInternalServerError)
	}
}

// metricsList := server.metrics.GetAll()

// var sb strings.Builder
// for key, value := range metricsList {
// 	sb.WriteString(fmt.Sprintf("%s: %v\n", key, value))
// }

// res.Write([]byte(sb.String()))
