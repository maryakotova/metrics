package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/maryakotova/metrics/internal/constants"
	"github.com/maryakotova/metrics/internal/models"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) PostgresStorage {
	return PostgresStorage{db: db}
}

func (ps PostgresStorage) Bootstrap(ctx context.Context) error {
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
		delta INTEGER DEFAULT 0
	);`

	tx.ExecContext(ctx, query)

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func (ps PostgresStorage) SetGauge(ctx context.Context, key string, value float64) (err error) {
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
	// _, err = ps.db.ExecContext(ctx, query, key, constants.Gauge, value)

	retries := 0
	for retries < 4 {
		_, err = ps.db.ExecContext(ctx, query, key, constants.Gauge, value)
		if err == nil {
			return nil
		}
		if !isRetriableError(err) {
			return err
		}
		retries++
		if retries == 4 {
			fmt.Println("ошибка соединения")
			return err
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}
	return err
}

func (ps PostgresStorage) SetCounter(ctx context.Context, key string, value *int64) (err error) {
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
	// _, err = ps.db.ExecContext(ctx, query, key, constants.Counter, value)

	retries := 0
	for retries < 4 {
		_, err = ps.db.ExecContext(ctx, query, key, constants.Counter, value)
		if err == nil {
			val, err := ps.GetCounter(ctx, key)
			if err == nil {
				*value = val
			}
			return nil
		}
		if !isRetriableError(err) {
			return err
		}
		retries++
		if retries == 4 {
			fmt.Println("ошибка соединения")
			return err
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	return err
}

func (ps PostgresStorage) GetGauge(ctx context.Context, key string) (value float64, err error) {
	query := `
	SELECT value FROM metrics WHERE id = $1 AND mtype = $2;
	`
	// row := ps.db.QueryRowContext(ctx, query, key, constants.Gauge)
	// if err = row.Scan(&value); err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		err = fmt.Errorf("значение метрики %s типа gauge не найдено", key)
	// 	} else {
	// 		err = fmt.Errorf("ошибка при сканировании значения: %v", err)
	// 	}
	// }

	retries := 0
	for retries < 4 {
		row := ps.db.QueryRowContext(ctx, query, key, constants.Gauge)
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
			fmt.Println("ошибка соединения")
			return
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	return
}

func (ps PostgresStorage) GetCounter(ctx context.Context, key string) (value int64, err error) {
	query := `
	SELECT delta FROM metrics WHERE id = $1 AND mtype = $2;
	`
	// row := ps.db.QueryRowContext(ctx, query, key, constants.Counter)
	// if err = row.Scan(&value); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		err = fmt.Errorf("значение метрики %s типа counter не найдено", key)
	// 	} else {
	// 		err = fmt.Errorf("ошибка при сканировании значения: %v", err)
	// 	}
	// }

	retries := 0
	for retries < 4 {
		row := ps.db.QueryRowContext(ctx, query, key, constants.Counter)
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
			fmt.Println("ошибка соединения")
			return
		}
		time.Sleep(time.Duration(retries*2+1) * time.Second) // Backoff: 1s, 3s, 5s
	}

	return
}

func (ps PostgresStorage) GetAllGauge(ctx context.Context) map[string]float64 {
	query := `
	SELECT id, value FROM metrics WHERE mtype = $1
	`
	// rows, err := ps.db.QueryContext(ctx, query, constants.Gauge)
	// if err != nil {
	// 	return nil
	// }

	var rows *sql.Rows
	var err error

	retries := 0
	for retries < 4 {
		rows, err = ps.db.QueryContext(ctx, query, constants.Gauge)
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil
		}
		retries++
		if retries == 4 {
			fmt.Println("ошибка соединения")
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

func (ps PostgresStorage) GetAllCounter(ctx context.Context) map[string]int64 {
	query := `
	SELECT id, delta FROM metrics WHERE mtype = $1
	`
	// rows, err := ps.db.QueryContext(ctx, query, constants.Counter)
	// if err != nil {
	// 	return nil
	// }

	var rows *sql.Rows
	var err error

	retries := 0
	for retries < 4 {
		rows, err = ps.db.QueryContext(ctx, query, constants.Counter)
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil
		}
		retries++
		if retries == 4 {
			fmt.Println("ошибка соединения")
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

func (ps PostgresStorage) GetAll(ctx context.Context) map[string]interface{} {

	query := `
	SELECT id, delta FROM metrics
	`
	// rows, err := ps.db.QueryContext(ctx, query)
	// if err != nil {
	// 	return nil
	// }
	var rows *sql.Rows
	var err error

	retries := 0
	for retries < 4 {
		rows, err = ps.db.QueryContext(ctx, query)
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil
		}
		retries++
		if retries == 4 {
			fmt.Println("ошибка соединения")
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

func (ps PostgresStorage) CheckConnection(ctx context.Context) (err error) {
	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ps.db.PingContext(context)
}

func (ps PostgresStorage) SaveMetrics(ctx context.Context, metrics []models.Metrics) (err error) {
	tx, err := ps.db.Begin()
	if err != nil {
		err = fmt.Errorf("failed to start transaction")
		return
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

	for _, metric := range metrics {
		if metric.ID == "" {
			err = fmt.Errorf("ошибка при сохранении в БД. Пустое имя метрики: %v", metric)
			return
		}
		var zero float64 = 0
		switch metric.MType {
		case constants.Gauge:
			if metric.Value == nil {
				metric.Value = &zero
			}
			_, err = tx.ExecContext(ctx, queryGauge, metric.ID, metric.MType, &metric.Value)
			if err != nil {
				return
			}
		case constants.Counter:
			_, err = tx.ExecContext(ctx, queryCounter, metric.ID, metric.MType, &metric.Delta)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("неверный формат для обновления метрик (недопустимый тип): %s", metric.MType)
			return
		}
	}

	return tx.Commit()
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
