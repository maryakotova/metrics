package agent

import (
	"sync"
	"time"
)

type Agent struct {
	ServerAddress  string
	ReportInterval int64
	PollInterval   int64
	SecretKey      string
	RateLimit      int
	retriesCount   int
	mutex          sync.Mutex
	metrics        map[string]interface{}
	WG             sync.WaitGroup
	pollTicker     *time.Ticker
	reportTicker   *time.Ticker
	sendQueue      chan Metrics
	resultQueue    chan Result
	publicKeyPath  string
	realIP         string
}

func New(cfg *Config) *Agent {
	return &Agent{
		ServerAddress:  cfg.ServerAddress,
		ReportInterval: cfg.ReportInterval,
		PollInterval:   cfg.PollInterval,
		SecretKey:      cfg.SecretKey,
		RateLimit:      cfg.RateLimit,
		retriesCount:   3,

		metrics:       make(map[string]interface{}),
		pollTicker:    time.NewTicker(time.Duration(cfg.PollInterval) * time.Second),
		reportTicker:  time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second),
		sendQueue:     make(chan Metrics, cfg.RateLimit),
		resultQueue:   make(chan Result, cfg.RateLimit),
		publicKeyPath: cfg.PublicCryptoKey,
		realIP:        cfg.RealIP,
	}
}
