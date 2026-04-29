package subscription

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env                 string        `env:"ENV"                            envDefault:"prod"`
	Host                string        `env:"HOST,notEmpty"`
	Port                string        `env:"PORT,notEmpty"`
	ReadHTO             time.Duration `env:"READ_HEADER_TIMEOUT"            envDefault:"5s"`
	ReadTO              time.Duration `env:"READ_TIMEOUT"                   envDefault:"15s"`
	WriteTO             time.Duration `env:"WRITE_TIMEOUT"                  envDefault:"30s"`
	IdleTO              time.Duration `env:"IDLE_TIMEOUT"                   envDefault:"120s"`
	RLRequestsPerSecond int           `env:"RATE_LIMIT_REQUESTS_PER_SECOND" envDefault:"10"`
	RLBurst             int           `env:"RATE_LIMIT_BURST"               envDefault:"30"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
