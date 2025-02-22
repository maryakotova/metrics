package config

import "github.com/maryakotova/metrics/internal/constants"

//-----------------------------------------------------------------------------------------------------------------------
// должны ли быть поля структуры Config публичными? или нужно сделать приватными и метод для доступа к каждому параметру?
//-----------------------------------------------------------------------------------------------------------------------

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	SecretKey string
}

type ServerConfig struct {
	ServerAddress   string
	StoreInterval   int64  // Интервал сохранения метрик на сервере в секундах
	FileStoragePath string // Имя файла, куда будут сохранены метрики
	Restore         bool   // Загружать или нет ранее сохраненные метрики из файла
}

type DatabaseConfig struct {
	DatabaseDsn string
	RetryCount  int //количество повторений при ошибках подключения
}

func NewConfig(flags *Flags) *Config {
	return &Config{
		Server: ServerConfig{
			ServerAddress:   flags.Server.ServerAddress,
			StoreInterval:   flags.Server.StoreInterval,
			FileStoragePath: flags.Server.FileStoragePath,
			Restore:         flags.Server.Restore,
		},
		Database: DatabaseConfig{
			DatabaseDsn: flags.Database.DatabaseDsn,
			RetryCount:  constants.RetryCount,
		},
		SecretKey: flags.SecretKey,
	}
}

func (cfg *Config) IsStoreInFileEnabled() bool {
	return cfg.Server.FileStoragePath != ""
}

func (cfg *Config) IsRestoreEnabled() bool {
	return cfg.Server.Restore
}

func (cfg *Config) IsDatabaseEnabled() bool {
	return cfg.Database.DatabaseDsn != ""
}

func (cfg *Config) GetStoreInterval() int64 {
	return cfg.Server.StoreInterval
}

func (cfg *Config) IsSyncStore() bool {
	return cfg.Server.StoreInterval == 0
}

func (cfg *Config) GetRetryCount() int {
	return cfg.Database.RetryCount
}
