package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/maryakotova/metrics/internal/filetransfer"
	"github.com/maryakotova/metrics/internal/storage"
)

var netAddress string
var interval int64
var filePath string
var restore bool

type Config struct {
	ServerAddress   string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func parseFlags() {

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&netAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&interval, "i", 0, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&filePath, "f", "./metricsStorage.json", "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&restore, "r", true, "Загрузка ранее сохранённые значения из указанного файла при старте сервера")
	flag.Parse()

	if cfg.ServerAddress != "" {
		netAddress = cfg.ServerAddress
	}
	if cfg.StoreInterval != 0 {
		interval = cfg.StoreInterval
	}
	if cfg.FileStoragePath != "" {
		filePath = cfg.FileStoragePath
	}
	if cfg.Restore {
		restore = cfg.Restore
	}
}

func UploadData(memStorage *storage.MemStorage) {
	if !restore {
		return
	}

	fileReader, err := filetransfer.NewFileReader(filePath)
	if err != nil {
		panic(err)
	}

	metrics, err := fileReader.ReadMetrics()
	if err != nil {
		return
	}

	defer fileReader.Close()

	if len(metrics) > 0 {
		for _, metric := range metrics {
			switch metric.MType {
			case "gauge":
				memStorage.SetGauge(metric.ID, *metric.Value)
			case "counter":
				memStorage.SetCounter(metric.ID, *metric.Delta)
			}
		}
	}
}
