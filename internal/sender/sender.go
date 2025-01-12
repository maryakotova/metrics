package sender

import (
	"fmt"
	"net/http"
)

func SendMetric(serverAddress string, metricType string, metricName string, metricValue interface{}) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", serverAddress, metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending metric:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Sent metric: %s/%s/%v\n", metricType, metricName, metricValue)
	return err
}
