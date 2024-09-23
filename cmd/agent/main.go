package main

import (
	"fmt"
	"log"

	"github.com/sudeeya/metrics-harvester/internal/agent"
	"github.com/sudeeya/metrics-harvester/internal/logging"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit)

	cfg, err := agent.NewConfig()
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

	logger.Info("Starting agent")
	agent := agent.NewAgent(logger, cfg)
	agent.Run()
}
