package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/agent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := agent.NewConfig()
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
	logger.Info("Starting agent")
	agent := agent.NewAgent(cfg, logger)
	agent.Run()
}
