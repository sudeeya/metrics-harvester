package server

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

const (
	defaultAddress         string = "localhost:8080"
	defaultCryptoKey       string = ""
	defaultDatabaseDSN     string = ""
	defaultKey             string = ""
	defaultLogLevel        string = "info"
	defaultStoreInterval   int64  = 300
	defaultFileStoragePath string = "metrics.json"
	defaultProfilerPort    int64  = 6060
	defaultRestore         bool   = true
)

type Config struct {
	Address         string `env:"ADDRESS"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	ProfilerPort    int64  `env:"PROFILER_PORT"`
	Restore         bool   `env:"RESTORE"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", defaultAddress, "Server IP address and port")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", defaultCryptoKey, "Path to the file where the RSA private key is saved")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaultDatabaseDSN, "Database DSN (e.g., user=postgres password=secret host=localhost port=5432 database=pgx_test sslmode=disable)")
	flag.StringVar(&cfg.Key, "k", defaultKey, "Key for HMAC hash")
	flag.StringVar(&cfg.LogLevel, "l", defaultLogLevel, "Log level: info, error, fatal")
	flag.Int64Var(&cfg.StoreInterval, "i", defaultStoreInterval, "The time interval in seconds after which metric values will be saved to the file")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultFileStoragePath, "Path to the file where the metric values are saved")
	flag.Int64Var(&cfg.ProfilerPort, "p", defaultProfilerPort, "The port on which pprof is running")
	flag.BoolVar(&cfg.Restore, "r", defaultRestore, "Determines whether previously saved values from a file will be loaded when the server starts")
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
