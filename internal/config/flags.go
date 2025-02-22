package config

import (
	"flag"
	"os"
	"strconv"
)

type Flags struct {
	Server struct {
		ServerAddress   string //`env:"ADDRESS"`
		StoreInterval   int64  //`env:"STORE_INTERVAL"`
		FileStoragePath string //`env:"FILE_STORAGE_PATH"`
		Restore         bool   //`env:"RESTORE"`
	}
	Database struct {
		DatabaseDsn string //`env:"DATABASE_DSN"`
	}
	SecretKey string //`env:"KEY"`
}

func ParseFlags() (*Flags, error) {

	var flags Flags
	var err error

	flag.StringVar(&flags.Server.ServerAddress, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&flags.Server.StoreInterval, "i", 0, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&flags.Server.FileStoragePath, "f", "./metricsStorage.json", "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&flags.Server.Restore, "r", true, "Загрузка ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&flags.Database.DatabaseDsn, "d", "", "Строка c адресом подключения к БД") //"host=localhost user=metrics password=test dbname=metrics sslmode=disable"
	flag.StringVar(&flags.SecretKey, "k", "", "Ключ для подписи передаваемых данных")

	//аргументы командной строки
	flag.Parse()

	// переменные окружения
	if envNetAddr := os.Getenv("ADDRESS"); envNetAddr != "" {
		flags.Server.ServerAddress = envNetAddr
	}

	if envInterval := os.Getenv("STORE_INTERVAL"); envInterval != "" {
		flags.Server.StoreInterval, err = strconv.ParseInt(envInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		flags.Server.FileStoragePath = envFilePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		flags.Server.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
	}

	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		flags.Database.DatabaseDsn = envDSN
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		flags.SecretKey = envKey
	}

	return &flags, nil

}
