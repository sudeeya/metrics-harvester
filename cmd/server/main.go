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
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Print(err)
		}
	}()
	memStorage := storage.NewMemStorage()
	logger.Info("Starting metrics-harvester")
	server := server.NewServer(logger, cfg, memStorage)
	server.Run()
}
