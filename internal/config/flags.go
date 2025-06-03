// В пакете config реализована возможность сервером принимать параметры конфигурации через флаги и переменные окружения.
// Приоритет параметров сервера:.
// - Если указана переменная окружения, то используется она.
// - Если нет переменной окружения, но есть флаг, то используется он.
// - Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"metrics/internal/constants"
	"metrics/internal/models"
	"os"
	"strconv"
)

// В структуру Flags сохраняются параметры конфигурации из флагов и переменных окружения.
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
	SecretKey        string //`env:"KEY"`
	PrivateCryptoKey string //`env:"CRYPTO_KEY"`
	ConfigPath       string //`env:"CONFIG"`
	ConfigPathShort  string
	TrustedSubnet    string //`env:"TRUSTED_SUBNET"`
}

// В функции ParseFlags происходит парсинг аргументов командной строки и получение значений из переменных окружения.
func ParseFlags() (*Flags, error) {

	var flags Flags
	var err error

	flag.StringVar(&flags.Server.ServerAddress, "a", constants.DefaultServerAddress, "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&flags.Server.StoreInterval, "i", constants.DefaultStoreInterval, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&flags.Server.FileStoragePath, "f", constants.DefaultStoreFile, "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&flags.Server.Restore, "r", constants.DefaultRestore, "Загрузка ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&flags.Database.DatabaseDsn, "d", "", "Строка c адресом подключения к БД") //"host=localhost user=metrics password=test dbname=metrics sslmode=disable"
	flag.StringVar(&flags.SecretKey, "k", "", "Ключ для подписи передаваемых данных")
	flag.StringVar(&flags.PrivateCryptoKey, "crypto-key", "", "Путь до файла с приватным ключом") //./key/private_key.pem
	flag.StringVar(&flags.ConfigPath, "config", "", "конфигурации сервера с помощью файла в формате JSON")
	flag.StringVar(&flags.ConfigPathShort, "c", "", "конфигурации сервера с помощью файла в формате JSON(shorthand)")
	flag.StringVar(&flags.TrustedSubnet, "t", "", "строковое представление бесклассовой адресации (CIDR) для проверки IP клиента")

	//аргументы командной строки
	flag.Parse()

	var serverConfig models.JSONConfigServer

	if flags.ConfigPath != "" || flags.ConfigPathShort != "" {
		var path string
		if flags.ConfigPath != "" {
			path = flags.ConfigPath
		} else {
			path = flags.ConfigPathShort
		}

		serverConfig, err = getJSONConfig(path)
		if err != nil {
			return nil, err
		}
	}

	// переменные окружения
	if envNetAddr := os.Getenv("ADDRESS"); envNetAddr != "" {
		flags.Server.ServerAddress = envNetAddr
	} else if flags.Server.ServerAddress == constants.DefaultServerAddress && serverConfig.Address != "" {
		flags.Server.ServerAddress = serverConfig.Address
	}

	if envInterval := os.Getenv("STORE_INTERVAL"); envInterval != "" {
		flags.Server.StoreInterval, err = strconv.ParseInt(envInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	} else if flags.Server.StoreInterval == constants.DefaultStoreInterval && serverConfig.StoreInterval != "" {
		interval, err := strconv.ParseInt(serverConfig.StoreInterval, 10, 64)
		if err == nil {
			flags.Server.StoreInterval = interval
		}
	}

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		flags.Server.FileStoragePath = envFilePath
	} else if flags.Server.FileStoragePath == constants.DefaultStoreFile && serverConfig.StoreFile != "" {
		flags.Server.FileStoragePath = serverConfig.StoreFile
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		flags.Server.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
	} else if flags.Server.Restore == constants.DefaultRestore && serverConfig.Restore != flags.Server.Restore {
		flags.Server.Restore = serverConfig.Restore
	}

	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		flags.Database.DatabaseDsn = envDSN
	} else if flags.Database.DatabaseDsn == "" && serverConfig.DatabaseDSN != "" {
		flags.Database.DatabaseDsn = serverConfig.DatabaseDSN
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		flags.SecretKey = envKey
	}

	if privateKey := os.Getenv("CRYPTO_KEY"); privateKey != "" {
		flags.PrivateCryptoKey = privateKey
	} else if flags.PrivateCryptoKey == "" && serverConfig.CryptoKey != "" {
		flags.PrivateCryptoKey = serverConfig.CryptoKey
	}

	if trustedSubnet := os.Getenv("TRUSTED_SUBNETY"); trustedSubnet != "" {
		flags.TrustedSubnet = trustedSubnet
	} else if flags.TrustedSubnet == "" && serverConfig.TrustedSubnet != "" {
		flags.TrustedSubnet = serverConfig.TrustedSubnet
	}

	return &flags, nil

}

func getJSONConfig(path string) (config models.JSONConfigServer, err error) {
	if path == "" {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("ошибка при чтении файла: %w", err)
		return
	}

	jsonErr := json.Unmarshal(data, &config)
	if jsonErr != nil {
		err = fmt.Errorf("ошибка преобразования JSON: %w", jsonErr)
		return
	}

	return
}
