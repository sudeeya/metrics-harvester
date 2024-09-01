package agent

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	defaultAddress         string = "localhost:8080"
	defaultBackoffSchedule string = "1,3,5"
	defaultKey             string = ""
	defaultLogLevel        string = "info"
	defaultPollInterval    int64  = 2
	defaultRateLimit       int64  = 16
	defaultReportInterval  int64  = 10
)

type Config struct {
<<<<<<< HEAD
	Address         string `env:"ADDRESS"`
	BackoffSchedule string `env:"BACKOFF_SCHEDULE"`
	Key             string `env:"KEY"`
	LogLevel        string `env:"LOG_LEVEL"`
	PollInterval    int64  `env:"POLL_INTERVAL"`
	RateLimit       int64  `env:"RATE_LIMIT"`
	ReportInterval  int64  `env:"REPORT_INTERVAL"`
=======
	Address         string `env:"ADDRESS" envDefault:"localhost:8080"`
	BackoffSchedule string `env:"BACKOFF_SCHEDULE" envDefault:"1,3,5"`
	Key             string `env:"KEY" envDefault:""`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	PollInterval    int64  `env:"POLL_INTERVAL" envDefault:"2"`
	RateLimit       int64  `env:"RATE_LIMIT" envDefault:"16"`
	ReportInterval  int64  `env:"REPORT_INTERVAL" envDefault:"10"`
>>>>>>> 6d7a864 (Add goroutine for gopsutil metrics)
}

func NewConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", defaultAddress, "Server IP address and port")
	flag.StringVar(&cfg.BackoffSchedule, "b", defaultBackoffSchedule, "Backoff schedule in seconds separated by commas")
	flag.StringVar(&cfg.Key, "k", defaultKey, "Key for HMAC hash")
	flag.StringVar(&cfg.LogLevel, "ll", defaultLogLevel, "Log level: info, error, fatal")
	flag.Int64Var(&cfg.PollInterval, "p", defaultPollInterval, "Polling interval in seconds")
	flag.Int64Var(&cfg.RateLimit, "l", defaultRateLimit, "Limit of requests")
	flag.Int64Var(&cfg.ReportInterval, "r", defaultReportInterval, "Report interval in seconds")
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
<<<<<<< HEAD
=======
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server IP address and port")
	flag.StringVar(&cfg.BackoffSchedule, "b", cfg.BackoffSchedule, "Backoff schedule in seconds separated by commas")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "Key for HMAC hash")
	flag.StringVar(&cfg.LogLevel, "ll", cfg.LogLevel, "Log level: info, error, fatal")
	flag.Int64Var(&cfg.PollInterval, "p", cfg.PollInterval, "Polling interval in seconds")
	flag.Int64Var(&cfg.RateLimit, "l", cfg.RateLimit, "Limit of requests")
	flag.Int64Var(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report interval in seconds")
	flag.Parse()
>>>>>>> 6d7a864 (Add goroutine for gopsutil metrics)
	if !strings.HasPrefix(cfg.Address, "http://") {
		cfg.Address = "http://" + cfg.Address
	}
	return &cfg, nil
}
