package agent

type Config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	PollInterval   int64  `env:"POLL_INTERVAL" envDefault:"2"`
	ReportInterval int64  `env:"REPORT_INTERVAL" envDefault:"10"`
}
