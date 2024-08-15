package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/agent"
	"github.com/sudeeya/metrics-harvester/internal/logging"
)

func main() {
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
