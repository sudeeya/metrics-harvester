package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/logging"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/server"
)

func main() {
	cfg, err := server.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := logging.NewLogger(cfg.LogLevel)
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
