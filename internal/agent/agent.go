package agent

import (
	"time"
)

type Agent struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
	SecretKey      string
	RateLimit      int
	retriesCount   int

	metrics map[string]interface{}
	//wg           sync.WaitGroup
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
	sendQueue    chan Metrics
	resultQueue  chan Result
}

func New(cfg *Config) *Agent {
	return &Agent{
		ServerAddress:  cfg.ServerAddress,
		ReportInterval: time.Duration(cfg.ReportInterval) * time.Second,
		PollInterval:   time.Duration(cfg.PollInterval) * time.Second,
		SecretKey:      cfg.SecretKey,
		RateLimit:      cfg.RateLimit,
		retriesCount:   3,

		metrics:      make(map[string]interface{}),
		pollTicker:   time.NewTicker(time.Duration(cfg.PollInterval)),
		reportTicker: time.NewTicker(time.Duration(cfg.ReportInterval)),
		sendQueue:    make(chan Metrics, cfg.RateLimit),
		resultQueue:  make(chan Result, cfg.RateLimit),
	}
}
