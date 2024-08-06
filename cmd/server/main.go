package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := server.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		log.Fatal(err)
	}
	memStorage := storage.NewMemStorage()
	server := server.NewServer(cfg, logger, memStorage)
	server.Run()
}
