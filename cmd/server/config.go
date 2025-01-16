package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

var netAddress string

type Config struct {
	ServerAddress string `env:"ADDRESS"`
}

func parseFlags() {

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&netAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.Parse()

	if cfg.ServerAddress != "" {
		netAddress = cfg.ServerAddress
	}
}

// func parseFlags() {

// 	flag.StringVar(&netAddress, "a", "localhost:8080", "Адрес и порт для HTTP-сервера")
// 	flag.Parse()

// 	if envServAddr := os.Getenv("ADDRESS"); envServAddr != "" {
// 		netAddress = envServAddr
// 	}
// }
