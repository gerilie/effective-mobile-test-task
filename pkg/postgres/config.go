package postgres

import "github.com/caarlos0/env/v11"

type Config struct {
	Host     string `env:"HOST"              envDefault:"postgres"`
	Port     string `env:"PORT"              envDefault:"5432"`
	User     string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	DB       string `env:"DB,notEmpty"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
