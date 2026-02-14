package logger

import (
	"github.com/carlcortright/k8s-scheduler/internal/config"

	"go.uber.org/zap"
)

// Simple singleton logger implementation using zap to log messages to the console
var logger *zap.Logger

func InitLogger(cfg config.Config) {
	if (cfg.Env == "development") {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
}

func GetLogger() *zap.Logger {
	// If logger is not initialized, initialize it and default to production
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return logger
}