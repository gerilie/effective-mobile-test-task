package logger

import "github.com/caarlos0/env/v11"

type Config struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
