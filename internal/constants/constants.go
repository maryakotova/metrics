// В пакете constants хранятся глобальные константы.
package constants

import "time"

const (
	Gauge                        = "gauge"
	Counter                      = "counter"
	PollCount                    = "PollCount"
	RetryCount            int    = 3
	HeaderSig             string = "HashSHA256"
	DefaultServerAddress         = "localhost:8080"
	DefaultStoreInterval  int64  = 300
	DefaultRestore        bool   = true
	DefaultStoreFile             = "./metricsStorage.json"
	DefaultReportInterval int64  = 5
	DefaultPollInterval   int64  = 2
	ShutdownTimeout              = 5 * time.Second
)
