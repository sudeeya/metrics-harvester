package server

type Config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}
