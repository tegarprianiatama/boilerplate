package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// InitLogger initializes the global Zap logger based on environment (e.g., "development" or "production").
func InitLogger(env string) *zap.Logger {
	once.Do(func() {
		var config zap.Config
		if env == "production" {
			config = zap.NewProductionConfig()
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		var err error
		log, err = config.Build()
		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}
	})
	return log
}

// Log returns the global logger instance.
func Log() *zap.Logger {
	if log == nil {
		return InitLogger("development")
	}
	return log
}

// Sync flushes any buffered log entries.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
