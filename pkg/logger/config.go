package logger

import (
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

// Config represents logger configuration loaded from environment variables.
type Config struct {
	// Level defines the logging level (e.g., debug, info, warn, error).
	Level zap.AtomicLevel `env:"LEVEL" envDefault:"info"`
}

// LoadConfig parses environment variables into Config.
func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
