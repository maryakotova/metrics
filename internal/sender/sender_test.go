package sender

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMetric_SuccessGauge(t *testing.T) {
	// Создаем тестовый сервер, который будет имитировать успешный ответ на POST запрос
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверим, что запрос это POST
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Ответим успешным кодом и без тела
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Подготовим данные для запроса
	serverAddress := ts.URL[len("http://"):]
	metricType := "Test1"
	metricName := "gauge"
	metricValue := 45.67

	// Выполним вызов функции SendMetric
	err := SendMetric(serverAddress, metricType, metricName, metricValue)

	// Проверим, что ошибок нет
	assert.NoError(t, err)
}

func TestSendMetric_SuccessCounter(t *testing.T) {
	// Создаем тестовый сервер, который будет имитировать успешный ответ на POST запрос
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверим, что запрос это POST
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Ответим успешным кодом и без тела
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Подготовим данные для запроса
	serverAddress := ts.URL[len("http://"):]
	metricType := "Test2"
	metricName := "counter"
	metricValue := "asd"

	// Выполним вызов функции SendMetric
	err := SendMetric(serverAddress, metricType, metricName, metricValue)

	// Проверим, что ошибок нет
	assert.NoError(t, err)
}

//не подходит, так как нужно проверять ответ сервера
// func TestSendMetric(t *testing.T) {
// 	type args struct {
// 		serverAddress string
// 		metricType    string
// 		metricName    string
// 		metricValue   interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "test gauge",
// 			args: args{
// 				serverAddress: "localhost:8080",
// 				metricType:    "gauge",
// 				metricName:    "Test1",
// 				metricValue:   123.456,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test gauge with -",
// 			args: args{
// 				serverAddress: "localhost:8080",
// 				metricType:    "gauge",
// 				metricName:    "Test2",
// 				metricValue:   -123.456,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test gauge with error",
// 			args: args{
// 				serverAddress: "localhost:8080",
// 				metricType:    "gauge",
// 				metricName:    "Test3",
// 				metricValue:   "abc",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test counter",
// 			args: args{
// 				serverAddress: "localhost:8080",
// 				metricType:    "counter",
// 				metricName:    "Test4",
// 				metricValue:   33,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test counter with -",
// 			args: args{
// 				serverAddress: "localhost:8080",
// 				metricType:    "counter",
// 				metricName:    "Test4",
// 				metricValue:   -10,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test counter with error",
// 			args: args{
// 				serverAddress: "localhost:8080",
// 				metricType:    "counter",
// 				metricName:    "Test4",
// 				metricValue:   321.456,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := SendMetric(tt.args.serverAddress, tt.args.metricType, tt.args.metricName, tt.args.metricValue); (err != nil) != tt.wantErr {
// 				t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
