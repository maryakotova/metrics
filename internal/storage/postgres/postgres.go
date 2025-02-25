package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/maryakotova/metrics/internal/config"
	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/models"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db     *sql.DB
	config *config.Config
	logger *zap.Logger
	m      sync.RWMutex
}

//------------------------------------------------------------------------------------
// нужно ли при таком создании конекшен вызывать defer db.Close(), если да, то где?
//------------------------------------------------------------------------------------

func NewPostgresStorage(cfg *config.Config, logger *zap.Logger) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", cfg.Database.DatabaseDsn)
	if err != nil {
		// можно ли тут вызывать панику? ----------------------------------------
		err = fmt.Errorf("не удалось подключиться к бд: %w", err)
		logger.Error(err.Error())
		return nil, err
	}
	return &PostgresStorage{
		db:     db,
		config: cfg,
		logger: logger,
	}, nil
}

func (ps *PostgresStorage) Bootstrap(ctx context.Context) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(50) PRIMARY KEY,
		mtype VARCHAR(10) NOT NULL,
		value DOUBLE PRECISION DEFAULT 0,
		delta BIGINT DEFAULT 0
	);`

	tx.ExecContext(ctx, query)

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func (ps *PostgresStorage) SetGauge(ctx context.Context, key string, value float64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}
	query := `
	INSERT INTO metrics (id, mtype, value)
	VALUES ($1, $2, $3) 
	ON CONFLICT (id) DO UPDATE
	SET mtype = EXCLUDED.mtype, value = EXCLUDED.value;
	`

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		_, err = ps.db.ExecContext(ctx, query, key, constants.Gauge, value)
		ps.m.Unlock()
		if err == nil {
			return nil
		}
		if !isRetriableError(err) {
			return err
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при сохранении counter в бд: %s, %v, %w", key, value, err)
			ps.logger.Error(err.Error())
			return err
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}
	return err
}

func (ps *PostgresStorage) SetCounter(ctx context.Context, key string, value *int64) (err error) {
	if key == "" {
		err = fmt.Errorf("имя метрики обязательно для заполнения")
		return
	}

	query := `
	INSERT INTO metrics (id, mtype, delta)
	VALUES ($1, $2, $3) 
	ON CONFLICT (id) DO UPDATE
	SET mtype = EXCLUDED.mtype, delta = metrics.delta + EXCLUDED.delta;
	`

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		_, err = ps.db.ExecContext(ctx, query, key, constants.Counter, value)
		ps.m.Unlock()
		if err == nil {
			val, err := ps.GetCounter(ctx, key)
			if err == nil {
				*value = val
			}
			return nil
		}
		if !isRetriableError(err) {
			err = fmt.Errorf("ошибка при сохранении counter в бд: %s,%v, %v, %w", key, value, *value, err)
			ps.logger.Error(err.Error())
			return err
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при сохранении counter в бд: %s,%v, %v, %w", key, value, *value, err)
			ps.logger.Error(err.Error())
			return err
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	return err
}

func (ps *PostgresStorage) GetGauge(ctx context.Context, key string) (value float64, err error) {
	query := `
	SELECT value FROM metrics WHERE id = $1 AND mtype = $2;
	`

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		row := ps.db.QueryRowContext(ctx, query, key, constants.Gauge)
		ps.m.Unlock()
		err = row.Scan(&value)
		if err == nil {
			return
		}
		if !isRetriableError(err) {
			if errors.Is(err, sql.ErrNoRows) {
				err = fmt.Errorf("значение метрики %s типа gauge не найдено (%w)", key, err)
				return
			}
			err = fmt.Errorf("ошибка при сканировании значения: %w", err)
			return
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при чтении gauge из бд: %s, %w", key, err)
			ps.logger.Error(err.Error())
			return
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	return
}

func (ps *PostgresStorage) GetCounter(ctx context.Context, key string) (value int64, err error) {
	query := `
	SELECT delta FROM metrics WHERE id = $1 AND mtype = $2;
	`

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		row := ps.db.QueryRowContext(ctx, query, key, constants.Counter)
		ps.m.Unlock()
		err = row.Scan(&value)
		if err == nil {
			return
		}
		if !isRetriableError(err) {
			if errors.Is(err, sql.ErrNoRows) {
				err = fmt.Errorf("значение метрики %s типа counter не найдено (%w)", key, err)
				return
			}
			err = fmt.Errorf("ошибка при сканировании значения: %w", err)
			return
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при чтении counter из бд: %s, %w", key, err)
			ps.logger.Error(err.Error())
			return
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	return
}

func (ps *PostgresStorage) GetAllGauge(ctx context.Context) map[string]float64 {
	query := `
	SELECT id, value FROM metrics WHERE mtype = $1
	`

	var rows *sql.Rows
	var err error

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		rows, err = ps.db.QueryContext(ctx, query, constants.Gauge)
		ps.m.Unlock()
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при чтении метрик из бд: %w", err)
			ps.logger.Error(err.Error())
			return nil
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	defer rows.Close()

	gaugeMetrics := make(map[string]float64)

	for rows.Next() {
		var id string
		var value float64

		if err := rows.Scan(&id, &value); err != nil {
			return nil
		}
		gaugeMetrics[id] = value
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return gaugeMetrics
}

func (ps *PostgresStorage) GetAllCounter(ctx context.Context) map[string]int64 {
	query := `
	SELECT id, delta FROM metrics WHERE mtype = $1
	`

	var rows *sql.Rows
	var err error

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		rows, err = ps.db.QueryContext(ctx, query, constants.Counter)
		ps.m.Unlock()
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при чтении метрик counter из бд: %w", err)
			ps.logger.Error(err.Error())
			return nil
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	defer rows.Close()

	counterMetrics := make(map[string]int64)
	for rows.Next() {
		var id string
		var value int64

		if err := rows.Scan(&id, &value); err != nil {
			return nil
		}
		counterMetrics[id] = value
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return counterMetrics
}

func (ps *PostgresStorage) GetAll(ctx context.Context) map[string]interface{} {

	query := `
	SELECT id, delta FROM metrics
	`
	var rows *sql.Rows
	var err error

	retries := 0
	for retries < 4 {
		ps.m.Lock()
		rows, err = ps.db.QueryContext(ctx, query)
		ps.m.Unlock()
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil
		}
		retries++
		if retries == 4 {
			err = fmt.Errorf("ошибка при чтении метрик из бд: %w", err)
			ps.logger.Error(err.Error())
			return nil
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}
	defer rows.Close()

	metrics := make(map[string]interface{})
	for rows.Next() {
		var id string
		var value int64

		if err := rows.Scan(&id, &value); err != nil {
			return nil
		}
		metrics[id] = value
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return metrics
}

func (ps *PostgresStorage) CheckConnection(ctx context.Context) (err error) {
	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ps.db.PingContext(context)
}

func (ps *PostgresStorage) SaveMetrics(ctx context.Context, metrics []models.Metrics) error {
	tx, err := ps.db.Begin()
	if err != nil {
		err = fmt.Errorf("failed to start transaction")
		return err
	}
	defer tx.Rollback()

	queryGauge := `
	INSERT INTO metrics (id, mtype, value)
	VALUES ($1, $2, $3) 
	ON CONFLICT (id) DO UPDATE
	SET mtype = EXCLUDED.mtype, value = EXCLUDED.value;
	`
	queryCounter := `
	INSERT INTO metrics (id, mtype, delta)
	VALUES ($1, $2, $3) 
	ON CONFLICT (id) DO UPDATE
	SET mtype = EXCLUDED.mtype, delta = metrics.delta + EXCLUDED.delta;
	`
	for i := 0; i <= ps.config.GetRetryCount(); i++ {

		var error error

		for _, metric := range metrics {
			if metric.ID == "" {
				error = fmt.Errorf("ошибка при сохранении в БД. Пустое имя метрики: %v", metric)
				break
			}
			var zero float64 = 0
			switch metric.MType {
			case constants.Gauge:
				if metric.Value == nil {
					metric.Value = &zero
				}
				ps.m.Lock()
				_, error = tx.ExecContext(ctx, queryGauge, metric.ID, metric.MType, &metric.Value)
				ps.m.Unlock()
				if error != nil {
					error = fmt.Errorf("ошибка при сохранении %s в бд: %s, %v, %w", metric.MType, metric.ID, &metric.Value, error)
				}
			case constants.Counter:
				ps.m.Lock()
				_, error = tx.ExecContext(ctx, queryCounter, metric.ID, metric.MType, &metric.Delta)
				ps.m.Unlock()
				if error != nil {
					error = fmt.Errorf("ошибка при сохранении %s в бд: %s, %v, %v, %w", metric.MType, metric.ID, &metric.Delta, metric.Delta, error)
				}
			default:
				error = fmt.Errorf("неверный формат для обновления метрик (недопустимый тип): %s", metric.MType)
			}
			if error != nil {
				break
			}
		}
		if error == nil {
			break
		}
		if !isRetriableError(err) {
			ps.logger.Error(error.Error())
			return error
		}
		if i == ps.config.GetRetryCount() {
			ps.logger.Error(error.Error())
			return error
		}
		time.Sleep(time.Duration(i*2+1) * time.Second) // Backoff: 1s, 3s, 5s

	}
	return tx.Commit()
}

func (ps *PostgresStorage) GetAllMetricsInJSON() []models.Metrics {
	metrics := []models.Metrics{}
	return metrics
}

func isRetriableError(err error) bool {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.Code == pgerrcode.ConnectionException ||
			pgerr.Code == pgerrcode.ConnectionDoesNotExist ||
			pgerr.Code == pgerrcode.ConnectionFailure ||
			pgerr.Code == pgerrcode.SQLClientUnableToEstablishSQLConnection ||
			pgerr.Code == pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection ||
			pgerr.Code == pgerrcode.TransactionResolutionUnknown ||
			pgerr.Code == pgerrcode.ProtocolViolation {
			return true
		}
	}
	return false
}
