package agent

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress  string
	ReportInterval int64
	PollInterval   int64
	SecretKey      string
	RateLimit      int
}

func ParseFlags() (*Config, error) {

	var err error
	var cfg Config

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&cfg.SecretKey, "k", "", "Ключ для подписи передаваемых данных")
	flag.IntVar(&cfg.RateLimit, "l", , "Количество одновременно исходящих запросов на сервер")

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

	return &cfg, nil

}
