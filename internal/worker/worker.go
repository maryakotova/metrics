// В пакете worker реализован вызов функции по таймеру и при завершении работы сервиса.
package worker

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TriggerGoFunc(ticker *time.Ticker, task func()) {

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				task()
			case <-signalChannel:
				task()
				os.Exit(0)
			}
		}
	}(ticker)

}
