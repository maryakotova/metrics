package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maryakotova/metrics/internal/authsign"
	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/models"
)

func SendMetric(serverAddress string, metricType string, metricName string, metricValue interface{}) error {

	metricForSend := models.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	var floatValue float64
	var intValue int64

	switch metricType {
	case constants.Gauge:
		switch v := metricValue.(type) {
		case int:
			floatValue = float64(v)
		case int8:
			floatValue = float64(v)
		case int16:
			floatValue = float64(v)
		case int32:
			floatValue = float64(v)
		case int64:
			floatValue = float64(v)
		case uint:
			floatValue = float64(v)
		case uint8:
			floatValue = float64(v)
		case uint16:
			floatValue = float64(v)
		case uint32:
			floatValue = float64(v)
		case uint64:
			floatValue = float64(v)
		case float32:
			floatValue = float64(v)
		case float64:
			floatValue = v
		}

	case constants.Counter:
		switch v := metricValue.(type) {
		case int64:
			intValue = v
		}
	}

	metricForSend.Delta = &intValue
	metricForSend.Value = &floatValue

	jsonData, err := json.Marshal(metricForSend)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n:", err)
		return err
	}

	url := fmt.Sprintf("http://%s/update/", serverAddress)

	buf := bytes.NewBuffer(nil)
	gzip := gzip.NewWriter(buf)
	_, err = gzip.Write([]byte(jsonData))
	if err != nil {
		fmt.Printf("Failed write data to compress temporary buffer: %v\n", err)
		return err
	}
	err = gzip.Close()
	if err != nil {
		fmt.Printf("Failed compress data: %v", err)
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return err
	}

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	// switch metricType {
	// case "gauge":
	// 	fmt.Printf("Sent metric: %s/%s/%v -/%v, response status: %s\n", metricType, metricName, floatValue, *metricForSend.Value, resp.Status)
	// case "counter":
	// 	fmt.Printf("Sent metric: %s/%s/%v -/%v, response status: %s\n", metricType, metricName, intValue, *metricForSend.Delta, resp.Status)
	// }

	return nil

}

func SendMetrics(serverAddress string, metrics map[string]interface{}) (err error) {
	for metricName, metricValue := range metrics {
		var metricType string
		if metricName == constants.PollCount {
			metricType = constants.Counter
		} else {
			metricType = constants.Gauge
		}
		err := SendMetric(serverAddress, metricType, metricName, metricValue)

		if err != nil {
			addText := fmt.Sprintf("Ошибка при отправке метрики %s\n", metricName)
			return fmt.Errorf("%s %w", addText, err)
		}
	}
	return nil
}

func SendMetricsBatch(serverAddress string, secretKey string, metrics []models.MetricsForSend) (err error) {
	if len(metrics) == 0 {
		err = fmt.Errorf("metrics table is empty")
		return
	}

	jsonData, err := json.Marshal(&metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	var hash string
	if secretKey != "" {
		hash = authsign.CalculateHash(jsonData, []byte(secretKey))
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to gzip writer: %v", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %v", err)
	}

	url := fmt.Sprintf("http://%s/updates/", serverAddress)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return err
	}

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "")
	if hash != "" {
		request.Header.Set(constants.HeaderSig, hash)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}
