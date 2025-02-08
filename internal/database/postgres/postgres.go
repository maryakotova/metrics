package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) PostgresStorage {
	return PostgresStorage{db: db}
}

func (ps PostgresStorage) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(50) PRIMARY KEY,
		mtype VARCHAR(10) NOT NULL,
		value DOUBLE PRECISION,
		delta INTEGER
	);`

	_, err := ps.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}
	return nil
}

func (ps PostgresStorage) SetGauge(key string, value float64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}
	//
	return
}

func (ps PostgresStorage) SetCounter(key string, value int64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}
	//
	return
}

func (ps PostgresStorage) GetGauge(key string) (value float64, err error) {

	//
	return
}

func (ps PostgresStorage) GetCounter(key string) (value int64, err error) {

	//
	return
}

func (ps PostgresStorage) GetAllGauge() map[string]float64 {
	return make(map[string]float64)
}

func (ps PostgresStorage) GetAllCounter() map[string]int64 {
	return make(map[string]int64)
}

func (ps PostgresStorage) GetAll() map[string]interface{} {

	return make(map[string]interface{})
}

func (ps PostgresStorage) CheckConnection(ctx context.Context) (err error) {
	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ps.db.PingContext(context)

}
