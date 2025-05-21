package agent

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress   string
	ReportInterval  int64
	PollInterval    int64
	SecretKey       string
	RateLimit       int
	PublicCryptoKey string //`env:"CRYPTO_KEY"`
}

func ParseFlags() (*Config, error) {

	var err error
	var cfg Config

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&cfg.ReportInterval, "r", 5, "Частота отправки метрик на сервер")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&cfg.SecretKey, "k", "", "Ключ для подписи передаваемых данных")
	flag.IntVar(&cfg.RateLimit, "l", 4, "Количество одновременно исходящих запросов на сервер")
	flag.StringVar(&cfg.PublicCryptoKey, "crypto-key", "./key/cert.pem", "Путь до файла с публичным ключом")

	//аргументы командной строки
	flag.Parse()

	// переменные окружения
	if envNetAddr := os.Getenv("ADDRESS"); envNetAddr != "" {
		cfg.ServerAddress = envNetAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		cfg.ReportInterval, err = strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		cfg.PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.SecretKey = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		cfg.RateLimit, err = strconv.Atoi(envRateLimit)
		if err != nil {
			return nil, err
		}
	}

	if publicKey := os.Getenv("CRYPTO_KEY"); publicKey != "" {
		cfg.PublicCryptoKey = publicKey
	}

	return &cfg, nil

}
