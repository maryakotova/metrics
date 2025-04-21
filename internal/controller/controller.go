package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"metrics/internal/constants"
	"metrics/internal/models"

	"go.uber.org/zap"
)

// const tplPath string = "./templates/metrics.html"
// const tplPath string = "templates/metrics.html"

type DataStorage interface {
	SetGauge(ctx context.Context, key string, value float64) (err error)
	SetCounter(ctx context.Context, key string, value *int64) (err error)
	SaveMetrics(ctx context.Context, metrics []models.Metrics) (err error)
	GetAllGauge(ctx context.Context) map[string]float64
	GetAllCounter(ctx context.Context) map[string]int64
	GetGauge(ctx context.Context, key string) (value float64, err error)
	GetCounter(ctx context.Context, key string) (value int64, err error)
	CheckConnection(ctx context.Context) (err error)
}

type Controller struct {
	storage DataStorage
	logger  *zap.Logger
}

func NewController(storage DataStorage, logger *zap.Logger) *Controller {
	return &Controller{
		storage: storage,
		logger:  logger,
	}
}

// ----------------------------------------------------------------------
//корректноли в качестве параметра из контроллера возвращать statusCode ?
// ----------------------------------------------------------------------

func (c *Controller) UpdateMetric(ctx context.Context, metric models.Metrics) (statusCode int, err error) {

	if metric.ID == "" {
		err = fmt.Errorf("ошибка при обновлении (имя метрики не заполнено)")
		c.logger.Error(err.Error())
		return http.StatusBadRequest, err
	}

	switch metric.MType {
	case constants.Gauge:
		err = c.storage.SetGauge(ctx, metric.ID, *metric.Value)
		if err != nil {
			err = fmt.Errorf("ошибка при обновлении %s типа %s: %w)", metric.ID, metric.MType, err)
			c.logger.Error(err.Error())
			statusCode = http.StatusInternalServerError
			return statusCode, err
		}
	case constants.Counter:
		err = c.storage.SetCounter(ctx, metric.ID, metric.Delta)
		if err != nil {
			err = fmt.Errorf("ошибка при обновлении %s типа %s: %w)", metric.ID, metric.MType, err)
			c.logger.Error(err.Error())
			statusCode = http.StatusInternalServerError
			return statusCode, err
		}
	default:
		err = fmt.Errorf("ошибка при обновлении %s: тип %s не поддерживается)", metric.ID, metric.MType)
		c.logger.Error(err.Error())
		statusCode = http.StatusBadRequest
		return statusCode, err
	}
	return
}

func (c *Controller) UpdateMetricFromString(ctx context.Context, mtype string, mname string, mvalue *string) (statusCode int, err error) {

	if mname == "" {
		err = fmt.Errorf("ошибка при обновлении (имя метрики не заполнено)")
		c.logger.Error(err.Error())
		return http.StatusBadRequest, err
	}

	switch mtype {
	case constants.Gauge:
		value, err := strconv.ParseFloat(*mvalue, 64)
		if err != nil {
			err = fmt.Errorf("ошибка при обновлении %s типа %s: неверный формат значения)", mname, mtype)
			c.logger.Error(err.Error())
			statusCode = http.StatusBadRequest
			return statusCode, err
		}
		err = c.storage.SetGauge(ctx, mname, value)
		if err != nil {
			err = fmt.Errorf("ошибка при обновлении %s типа %s: %w)", mname, mtype, err)
			c.logger.Error(err.Error())
			statusCode = http.StatusInternalServerError
			return statusCode, err
		}
	case constants.Counter:
		value, err := strconv.ParseInt(*mvalue, 10, 64)
		if err != nil {
			err = fmt.Errorf("ошибка при обновлении %s типа %s: неверный формат значения)", mname, mtype)
			c.logger.Error(err.Error())
			statusCode = http.StatusBadRequest
			return statusCode, err
		}
		err = c.storage.SetCounter(ctx, mname, &value)
		if err != nil {
			err = fmt.Errorf("ошибка при обновлении %s типа %s: %w)", mname, mtype, err)
			c.logger.Error(err.Error())
			statusCode = http.StatusInternalServerError
			return statusCode, err
		}
		*mvalue = strconv.FormatInt(value, 10)
	default:
		err = fmt.Errorf("ошибка при обновлении %s: тип %s не поддерживается)", mname, mtype)
		c.logger.Error(err.Error())
		statusCode = http.StatusBadRequest
		return statusCode, err
	}
	return
}

func (c *Controller) GetOneMetric(ctx context.Context, metric *models.Metrics) (statusCode int, err error) {
	switch metric.MType {
	case constants.Gauge:
		value, err := c.storage.GetGauge(ctx, metric.ID)
		if err != nil {
			err = fmt.Errorf("не удалось получить данные для метрики %s типа %s: %w", metric.ID, metric.MType, err)
			c.logger.Error(err.Error())
			statusCode = http.StatusNotFound
			return statusCode, err
		}
		metric.Value = &value
	case constants.Counter:
		delta, err := c.storage.GetCounter(ctx, metric.ID)

		if err != nil {
			err = fmt.Errorf("не удалось получить данные для метрики %s типа %s: %w", metric.ID, metric.MType, err)
			c.logger.Error(err.Error())
			statusCode = http.StatusNotFound
			return statusCode, err
		}
		metric.Delta = &delta
	case "":
		err = fmt.Errorf("ошибка при получении метрики %s: тип обязателем для заполнения)", metric.ID)
		c.logger.Error(err.Error())
		statusCode = http.StatusBadRequest
		return statusCode, err
	default:
		err = fmt.Errorf("ошибка при получении метрики %s: тип %s не поддерживается)", metric.ID, metric.MType)
		c.logger.Error(err.Error())
		statusCode = http.StatusBadRequest
		return statusCode, err
	}
	return
}

func (c *Controller) GetAllGauge(ctx context.Context) map[string]float64 {
	return c.storage.GetAllGauge(ctx)
}
func (c *Controller) GetAllCounter(ctx context.Context) map[string]int64 {
	return c.storage.GetAllCounter(ctx)
}

func (c *Controller) CheckConnection(ctx context.Context) (err error) {
	err = c.storage.CheckConnection(ctx)
	if err != nil {
		c.logger.Error(err.Error())
	}
	return
}

func (c *Controller) SaveMetrics(ctx context.Context, metrics []models.Metrics) (statusCode int, err error) {
	for _, metric := range metrics {
		if metric.ID == "" {
			err = fmt.Errorf("ошибка при обновлении (имя метрики типа %s не заполнено)", metric.MType)
			c.logger.Error(err.Error())
			statusCode = http.StatusBadRequest
			return statusCode, err
		}
		switch metric.MType {
		case "":
			err = fmt.Errorf("ошибка при обновлении (тип метрики %s не заполнен)", metric.ID)
			c.logger.Error(err.Error())
			statusCode = http.StatusBadRequest
			return statusCode, err
		case constants.Counter:
		case constants.Gauge:
		default:
			err = fmt.Errorf("ошибка при получении метрики %s: тип %s не поддерживается)", metric.ID, metric.MType)
			c.logger.Error(err.Error())
			statusCode = http.StatusBadRequest
			return statusCode, err
		}
	}
	err = c.storage.SaveMetrics(ctx, metrics)
	if err != nil {
		c.logger.Error(err.Error())
		statusCode = http.StatusInternalServerError
		return statusCode, err
	}
	return
}
