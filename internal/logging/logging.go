// Package logging provides a way to create logger.
package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log levels
const (
	Info  = "info"
	Error = "error"
	Fatal = "fatal"
)

func NewLogger(logLevel string) (*zap.Logger, error) {
	logConfig := zap.NewDevelopmentConfig()
	switch logLevel {
	case Info:
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case Error:
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case Fatal:
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		log.Fatalf("unknown log level: %s", logLevel)
	}
	return logConfig.Build()
}
