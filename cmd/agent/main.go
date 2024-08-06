package main

import (
	"log"

	"github.com/sudeeya/metrics-harvester/internal/agent"
)

func main() {
	cfg, err := agent.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	agent := agent.NewAgent(cfg)
	agent.Run()
}
