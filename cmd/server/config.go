package main

import (
	"flag"
	"os"
)

var netAddress string

func parseFlags() {

	flag.StringVar(&netAddress, "a", "localhost:8080", "Адрес и порт для HTTP-сервера")
	flag.Parse()

	if envServAddr := os.Getenv("ADDRESS"); envServAddr != "" {
		netAddress = envServAddr
	}
}
