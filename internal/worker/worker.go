package worker

import "time"

func InitPeriodicFunc(interval int64, task func()) {

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		task()
	}
}
