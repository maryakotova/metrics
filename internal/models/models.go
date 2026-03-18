// В пакете models хранятся структуры для работы с данными.
package models

// Структура используется сервером для получения данных от агента
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Структура используется агентом для сбора и отправки данных на сервер
type MetricsForSend struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Конфигурации сервера с помощью файла в формате JSON
type JSONConfigServer struct {
	Address       string `json:"address"`        // аналог переменной окружения ADDRESS или флага -a
	Restore       bool   `json:"restore"`        // аналог переменной окружения RESTORE или флага -r
	StoreInterval string `json:"store_interval"` // аналог переменной окружения STORE_INTERVAL или флага -i
	StoreFile     string `json:"store_file"`     // аналог переменной окружения STORE_FILE или -f
	DatabaseDSN   string `json:"database_dsn"`   // аналог переменной окружения DATABASE_DSN или флага -d
	CryptoKey     string `json:"crypto_key"`     // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
	TrustedSubnet string `json:"trusted_subnet"` // аналог переменной окружения TRUSTED_SUBNET или флага -t
}

// Конфигурации агента с помощью файла в формате JSON
type JSONConfigAgent struct {
	Address        string `json:"address"`         // аналог переменной окружения ADDRESS или флага -a
	ReportInterval string `json:"report_interval"` // аналог переменной окружения REPORT_INTERVAL или флага -r
	PollInterval   string `json:"poll_interval"`   // аналог переменной окружения POLL_INTERVAL или флага -p
	CryptoKey      string `json:"crypto_key"`      // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}
