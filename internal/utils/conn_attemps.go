package utils

import (
	"time"

	"github.com/ellofae/go-concurrency-process/pkg/logger"
)

func ConnectionAttemps(conn_func func() error, attemps int, delay time.Duration) (err error) {
	logger := logger.GetLogger()

	for i := 0; i < attemps; i++ {
		err = conn_func()
		if err != nil {
			logger.Warn("Attempting to connect", "current attemp", i+1, "appemps left", attemps-i-1)
			time.Sleep(delay)
			continue
		}
	}
	return err
}
