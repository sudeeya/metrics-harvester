package server

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval   int64  `env:"STORE_INTERVAL" envDefault:"300"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	Restore         bool   `env:"RESTORE" envDefault:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server IP address and port")
	flag.Int64Var(&cfg.StoreInterval, "i", cfg.StoreInterval, "The time interval in seconds after which metric values will be saved to the file")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Path to the file where the metric values are saved")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "Determines whether previously saved values from a file will be loaded when the server starts")
	flag.Parse()
	return &cfg, nil
}
