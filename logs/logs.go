package logs

import (
	"log"
	"sync"
)

var (
	logger log.Logger
	once   sync.Once
)

func GetLogger() *log.Logger {
	once.Do(func() {
		logger = *log.Default()
	})
	return &logger
}
