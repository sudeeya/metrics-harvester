package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := server.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logConfig := zap.NewDevelopmentConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, err := logConfig.Build()
	defer logger.Sync()
	if err != nil {
		log.Fatal(err)
	}
	memStorage := storage.NewMemStorage()
	logger.Info("Starting metrics-harvester")
	server := server.NewServer(cfg, logger, memStorage)
	server.Run()
}
