package trustedip

import (
	"metrics/internal/config"
	"net"
	"net/http"

	"go.uber.org/zap"
)

type CheckIPMiddleware struct {
	config *config.Config
	logger *zap.Logger
	ipNet  *net.IPNet
}

func NewCheckIPMW(cfg *config.Config, logger *zap.Logger) *CheckIPMiddleware {
	ipmw := &CheckIPMiddleware{
		config: cfg,
		logger: logger,
	}

	if cfg.TrustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			//error
		} else {
			ipmw.ipNet = ipNet
		}
	}

	return ipmw
}

func (chip *CheckIPMiddleware) CheckIP(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if chip.ipNet != nil {
			// получить IP адрес из заголовка
			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
			if ip == nil {
				http.Error(w, "неверный IP адрес", http.StatusForbidden)
				return
			}

			// проверить входит ли полученный IP адрес в список доверенных
			if !chip.ipNet.Contains(ip) {
				http.Error(w, "неверный IP адрес", http.StatusForbidden)
				return
			}

		}

		h.ServeHTTP(ow, r)
	}
}
