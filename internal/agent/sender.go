package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"metrics/internal/authsign"
	"metrics/internal/constants"
	"metrics/internal/models"
)

func (agent *Agent) PrepareMetrics(metrics map[string]interface{}) []models.MetricsForSend {
	metricsForJSON := []models.MetricsForSend{}

	for key, value := range metrics {

		metric := models.MetricsForSend{ID: key}

		if key == constants.PollCount {
			metric.MType = constants.Counter
			switch v := value.(type) {
			case int64:
				metric.Delta = v
			case int:
				metric.Delta = int64(v)
			}

		} else {
			var floatValue float64
			metric.MType = constants.Gauge
			switch v := value.(type) {
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
			metric.Value = floatValue
		}

		metricsForJSON = append(metricsForJSON, metric)
	}
	return metricsForJSON
}

func (agent *Agent) SendMetricsBatch(metrics []models.MetricsForSend) error {
	var err error
	if len(metrics) == 0 {
		err = fmt.Errorf("metrics table is empty")
		return err
	}

	for i := 0; i <= agent.retriesCount; i++ {

		err = agent.doUpdatesRequest(metrics)
		if err == nil {
			break
		}

		if !isRetriableError(err) {
			return err
		}

		if i == agent.retriesCount {
			fmt.Println("ошибка соединения: %w", err)
			return err
		}

		time.Sleep(time.Duration(i*2+1) * time.Second)

	}
	return nil
}

func (agent *Agent) doUpdatesRequest(metrics []models.MetricsForSend) error {
	jsonData, err := json.Marshal(&metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	var hash string
	if agent.SecretKey != "" {
		hash = authsign.CalculateHash(jsonData, []byte(agent.SecretKey))
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to gzip writer: %w", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	url := fmt.Sprintf("http://%s/updates/", agent.ServerAddress)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		fmt.Println("Error sending request: %w\n", err)
		return err
	}

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "")
	if hash != "" {
		request.Header.Set(constants.HeaderSig, hash)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error sending request: %w\n", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

func isRetriableError(err error) bool {
	var opErr *net.OpError
	return errors.As(err, &opErr)
}
