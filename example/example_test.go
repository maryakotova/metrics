// Пакет example_test предоставляет примеры использования сервера метрик.
package example_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"metrics/internal/constants"
	"metrics/internal/models"
	"net/http"
	"time"
)

const ServerAddr = "http://localhost:8080"

func Example_updateJSON() {
	resp, err := sendMetric()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	// Output: Статус ответа: 200
}

func Example_getValue() {
	resp, err := sendMetric()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Запрос не был обработан успешноб статус: %d\n", resp.StatusCode)
		return
	}

	time.Sleep(1 * time.Second)

	metric := models.Metrics{
		ID:    "TestGauge",
		MType: constants.Gauge,
	}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		fmt.Printf("Ошибка при переводе в JSON: %v\n", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, ServerAddr+"/value", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}

	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Ошибка при отпарке запроса: %w\n", err)
		return
	}

	defer resp.Body.Close()

	var result models.Metrics
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Ошибка при переводе ответа из JSON: %v\n", err)
		return
	}

	fmt.Printf("Значение метрики %s: %f\n", result.ID, *result.Value)
	// Output: Значение метрики TestGauge: 111.223000

}

func sendMetric() (resp *http.Response, err error) {
	metric := models.Metrics{
		ID:    "TestGauge",
		MType: constants.Gauge,
		Value: new(float64),
	}
	*metric.Value = 111.223

	jsonData, err := json.Marshal(metric)
	if err != nil {
		err = fmt.Errorf("Ошибка при переводе в JSON: %v\n", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, ServerAddr+"/update", bytes.NewBuffer(jsonData))
	if err != nil {
		err = fmt.Errorf("Ошибка при создании запроса: %v\n", err)
		return nil, err
	}

	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		err = fmt.Errorf("Ошибка при отпарке запроса: %w\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
