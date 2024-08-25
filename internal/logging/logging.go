package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel string) (*zap.Logger, error) {
	logConfig := zap.NewDevelopmentConfig()
	switch logLevel {
	case "info":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "error":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		log.Fatalf("unknown log level: %s", logLevel)
	}
	return logConfig.Build()
}
