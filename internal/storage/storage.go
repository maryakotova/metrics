// В пакете storage реализован паттерн factory для создания необходимого типа хранилища.
package storage

import (
	"context"

	"metrics/internal/config"
	"metrics/internal/models"
	"metrics/internal/storage/inmemory"
	"metrics/internal/storage/postgres"

	"go.uber.org/zap"
)

// Интерфейс Storage необходим для взаимодействия с хранилищем метрик.
type Storage interface {
	SetGauge(ctx context.Context, key string, value float64) (err error)
	SetCounter(ctx context.Context, key string, value *int64) (err error)
	SaveMetrics(ctx context.Context, metrics []models.Metrics) (err error)
	GetAllGauge(ctx context.Context) map[string]float64
	GetAllCounter(ctx context.Context) map[string]int64
	GetGauge(ctx context.Context, key string) (value float64, err error)
	GetCounter(ctx context.Context, key string) (value int64, err error)
	CheckConnection(ctx context.Context) (err error)
	GetAllMetricsInJSON() []models.Metrics
}

type StorageFactory struct{}

// NewStorage реализует паттерн Factory и создает необходимый вид хранилища в зависимости от полученных настроек.
func (f *StorageFactory) NewStorage(cfg *config.Config, logger *zap.Logger) (Storage, error) {
	if cfg.IsDatabaseEnabled() {
		postgres, err := postgres.NewPostgresStorage(cfg, logger)
		if err != nil {
			return nil, err
		}
		err = postgres.Bootstrap(context.TODO())
		return postgres, err
	} else {
		inmemory := inmemory.NewMemStorage()
		if cfg.IsRestoreEnabled() {
			inmemory.UploadData(cfg.Server.FileStoragePath)
		}
		return inmemory, nil
	}
}
