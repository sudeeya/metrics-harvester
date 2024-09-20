package main

import (
	"fmt"
	"log"

	"github.com/sudeeya/metrics-harvester/internal/logging"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"github.com/sudeeya/metrics-harvester/internal/repository/database"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit)

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
		repository = database.NewDatabase(cfg.DatabaseDSN)
	default:
		repository = storage.NewMemStorage()
	}

	logger.Info("Starting metrics-harvester")
	server := server.NewServer(logger, cfg, repository)
	server.Run()
}
