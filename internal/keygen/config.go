package keygen

import "flag"

type Config struct {
	PrivateKeyPath string
	PublicKeyPath  string
}

func NewConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.PrivateKeyPath, "private", "private.pem", "Path to the file where the RSA private key will be saved")
	flag.StringVar(&cfg.PublicKeyPath, "public", "public.pem", "Path to the file where the RSA public key will be saved")
	flag.Parse()
	return &cfg, nil
}
