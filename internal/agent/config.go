package agent

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	LogLevel       string `env:"LOG_LEVEL" envDefault:"info"`
	PollInterval   int64  `env:"POLL_INTERVAL" envDefault:"2"`
	ReportInterval int64  `env:"REPORT_INTERVAL" envDefault:"10"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server IP address and port")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "Log level: info, error, fatal")
	flag.Int64Var(&cfg.PollInterval, "p", cfg.PollInterval, "Polling interval in seconds")
	flag.Int64Var(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report interval in seconds")
	flag.Parse()
	if !strings.HasPrefix(cfg.Address, "http://") {
		cfg.Address = "http://" + cfg.Address
	}
	return &cfg, nil
}
