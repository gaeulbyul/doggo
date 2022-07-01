package utils

import (
	"github.com/mr-karan/logf"
)

// InitLogger initializes logger.
func InitLogger() *logf.Logger {
	logger := logf.New()
	logger.SetColorOutput(true)
	return logger
}
