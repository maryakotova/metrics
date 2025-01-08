package main

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

var (
	serverAddress  string
	pollInterval   int64
	reportInterval int64
)

type Config struct {
	ServerAddress  string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
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
	flag.Parse()

	if cfg.ServerAddress != "" {
		serverAddress = cfg.ServerAddress
	}
	if cfg.ReportInterval != 0 {
		reportInterval = int64(cfg.ReportInterval.Seconds())
	}
	if cfg.PollInterval != 0 {
		pollInterval = int64(cfg.PollInterval.Seconds())
	}

}
