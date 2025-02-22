package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	SecretKey      string `env:"KEY"`
}

func parseFlags() {

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&serverAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&reportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.Int64Var(&pollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&secretKey, "k", "", "Ключ для подписи передаваемых данных")
	flag.Parse()

	if cfg.ServerAddress != "" {
		serverAddress = cfg.ServerAddress
	}
	if cfg.ReportInterval != 0 {
		reportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		pollInterval = cfg.PollInterval
	}

	if cfg.SecretKey != "" {
		secretKey = cfg.SecretKey
	}

}
