// В пакете constants хранятся глобальные константы.
package constants

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
)
