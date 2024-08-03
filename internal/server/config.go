package server

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server IP address and port")
	flag.Parse()
	return &cfg, nil
}
