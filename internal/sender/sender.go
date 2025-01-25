package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
	case "gauge":
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

	case "counter":
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	switch metricType {
	case "gauge":
		fmt.Printf("Sent metric: %s/%s/%v -/%v, response status: %s\n", metricType, metricName, floatValue, *metricForSend.Value, resp.Status)
	case "counter":
		fmt.Printf("Sent metric: %s/%s/%v -/%v, response status: %s\n", metricType, metricName, intValue, *metricForSend.Delta, resp.Status)
	}

	return nil

}
