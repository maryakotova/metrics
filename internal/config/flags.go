package config

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

// type Flags struct {
// 	Server struct {
// 		ServerAddress   string `env:"ADDRESS"`
// 		StoreInterval   int64  `env:"STORE_INTERVAL"`
// 		FileStoragePath string `env:"FILE_STORAGE_PATH"`
// 		Restore         bool   `env:"RESTORE"`
// 	}
// 	Database struct {
// 		DatabaseDsn string `env:"DATABASE_DSN"`
// 	}
// }

// func parseFlags() {

// 	var flags Flags

// 	// переменные окружения
// 	err := env.Parse(&flags)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	flag.StringVar(&netAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
// 	flag.Int64Var(&interval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
// 	flag.StringVar(&filePath, "f", "./metricsStorage.json", "Путь до файла, куда сохраняются текущие значения")
// 	flag.BoolVar(&restore, "r", true, "Загрузка ранее сохранённые значения из указанного файла при старте сервера")
// 	flag.StringVar(&dbDsn, "d", "", "Строка c адресом подключения к БД") //"host=localhost user=metrics password=test dbname=metrics sslmode=disable"

// 	//аргументы командной строки
// 	flag.Parse()

// 	if flags.Server.ServerAddress != "" {
// 		netAddress = flags.Server.ServerAddress
// 	}
// 	if flags.Server.StoreInterval != 0 {
// 		interval = flags.Server.StoreInterval
// 	}
// 	if flags.Server.FileStoragePath != "" {
// 		filePath = flags.Server.FileStoragePath
// 	}
// 	if flags.Server.Restore {
// 		restore = flags.Server.Restore
// 	}
// 	if flags.Database.DatabaseDsn != "" {
// 		dbDsn = flags.Database.DatabaseDsn
// 	}

// 	// если данные сохраняются в БД, путь к файлу для сохранения данных не требуется
// 	if dbDsn != "" {
// 		filePath = ""
// 	}
// }
