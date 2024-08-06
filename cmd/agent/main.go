package main

import (
	"github.com/sudeeya/metrics-harvester/internal/agent"
)

func main() {
	cfg, err := agent.NewConfig()
	if err != nil {
		panic(err)
	}
	agent := agent.NewAgent(cfg)
	agent.Run()
}
