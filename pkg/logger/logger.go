package logger

import (
	"go.uber.org/zap"
)

func New() *zap.Logger {
	// Production JSON logger
	log, _ := zap.NewProduction()
	return log
}
