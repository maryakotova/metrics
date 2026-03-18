package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"metrics/internal/constants"
	"metrics/internal/models"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress   string
	ReportInterval  int64
	PollInterval    int64
	SecretKey       string
	RateLimit       int
	PublicCryptoKey string //`env:"CRYPTO_KEY"`
	ConfigPath      string //`env:CONFIG`
	ConfigPathShort string
	// не знаю, как еще можно получить IP адрес агента, поэтому решила использовать переменные окружения и флаги
	RealIP string //`env:REAL_IP`
}

func ParseFlags() (*Config, error) {

	var err error
	var cfg Config

	flag.StringVar(&cfg.ServerAddress, "a", constants.DefaultServerAddress, "Адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&cfg.ReportInterval, "r", constants.DefaultReportInterval, "Частота отправки метрик на сервер")
	flag.Int64Var(&cfg.PollInterval, "p", constants.DefaultPollInterval, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&cfg.SecretKey, "k", "", "Ключ для подписи передаваемых данных")
	flag.IntVar(&cfg.RateLimit, "l", 4, "Количество одновременно исходящих запросов на сервер")
	flag.StringVar(&cfg.PublicCryptoKey, "crypto-key", "", "Путь до файла с публичным ключом") //./key/cert.pem
	flag.StringVar(&cfg.ConfigPath, "config", "", "конфигурации сервера с помощью файла в формате JSON")
	flag.StringVar(&cfg.ConfigPathShort, "c", "", "конфигурации сервера с помощью файла в формате JSON(shorthand)")
	flag.StringVar(&cfg.RealIP, "i", "", "IP адрес агента для заполнения заголовка X-Real-IP")

	//аргументы командной строки
	flag.Parse()

	var agentConfig models.JSONConfigAgent

	if cfg.ConfigPath != "" || cfg.ConfigPathShort != "" {
		var path string
		if cfg.ConfigPath != "" {
			path = cfg.ConfigPath
		} else {
			path = cfg.ConfigPathShort
		}

		agentConfig, err = getJSONConfig(path)
		if err != nil {
			return nil, err
		}
	}

	// переменные окружения
	if envNetAddr := os.Getenv("ADDRESS"); envNetAddr != "" {
		cfg.ServerAddress = envNetAddr
	} else if cfg.ServerAddress == constants.DefaultServerAddress && agentConfig.Address != "" {
		cfg.ServerAddress = agentConfig.Address
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		cfg.ReportInterval, err = strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	} else if cfg.ReportInterval == constants.DefaultReportInterval && agentConfig.ReportInterval != "" {
		reportInterval, err := strconv.ParseInt(agentConfig.ReportInterval, 10, 64)
		if err == nil {
			cfg.ReportInterval = reportInterval
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		cfg.PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	} else if cfg.PollInterval == constants.DefaultPollInterval && agentConfig.PollInterval != "" {
		pollInterval, err := strconv.ParseInt(agentConfig.PollInterval, 10, 64)
		if err == nil {
			cfg.PollInterval = pollInterval
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

	if publicKey := os.Getenv("CRYPTO_KEY"); publicKey != "" {
		cfg.PublicCryptoKey = publicKey
	} else if cfg.PublicCryptoKey != "" && agentConfig.CryptoKey != "" {
		cfg.PublicCryptoKey = agentConfig.CryptoKey
	}

	if realIP := os.Getenv("REAL_IP"); realIP != "" {
		cfg.RealIP = realIP
	}

	return &cfg, nil

}

func getJSONConfig(path string) (config models.JSONConfigAgent, err error) {
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
