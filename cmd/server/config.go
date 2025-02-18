package main

// import (
// 	"flag"
// 	"log"

// 	"github.com/caarlos0/env"
// )

// var netAddress string
// var interval int64
// var filePath string
// var restore bool
// var dbDsn string

// type Config struct {
// 	ServerAddress   string `env:"ADDRESS"`
// 	StoreInterval   int64  `env:"STORE_INTERVAL"`
// 	FileStoragePath string `env:"FILE_STORAGE_PATH"`
// 	Restore         bool   `env:"RESTORE"`
// 	DatabaseDsn     string `env:"DATABASE_DSN"`
// }

// // func NewConfig() *Config {
// //     return &Config{
// //         ServerAddress:   "localhost:8080",
// //         StoreInterval:   300,
// //         FileStoragePath: "./metricsStorage.json",
// //         Restore:         true,

// //     }
// // }

// func parseFlags() {

// 	var cfg Config
// 	err := env.Parse(&cfg)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	flag.StringVar(&netAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
// 	flag.Int64Var(&interval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
// 	flag.StringVar(&filePath, "f", "./metricsStorage.json", "Путь до файла, куда сохраняются текущие значения")
// 	flag.BoolVar(&restore, "r", true, "Загрузка ранее сохранённые значения из указанного файла при старте сервера")
// 	flag.StringVar(&dbDsn, "d", "", "Строка c адресом подключения к БД") //"host=localhost user=metrics password=test dbname=metrics sslmode=disable"

// 	flag.Parse()

// 	if cfg.ServerAddress != "" {
// 		netAddress = cfg.ServerAddress
// 	}
// 	if cfg.StoreInterval != 0 {
// 		interval = cfg.StoreInterval
// 	}
// 	if cfg.FileStoragePath != "" {
// 		filePath = cfg.FileStoragePath
// 	}
// 	if cfg.Restore {
// 		restore = cfg.Restore
// 	}
// 	if cfg.DatabaseDsn != "" {
// 		dbDsn = cfg.DatabaseDsn
// 	}

// 	if dbDsn != "" {
// 		filePath = ""
// 	}
// }
