package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/logging"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"github.com/sudeeya/metrics-harvester/internal/repository/database"
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
	var repository repo.Repository
	switch {
	case cfg.DatabaseDSN != "":
		repository, err = database.NewDatabase(cfg.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error())
		}
	default:
		repository = storage.NewMemStorage()
	}
	logger.Info("Starting metrics-harvester")
	server := server.NewServer(logger, cfg, repository)
	server.Run()
}
