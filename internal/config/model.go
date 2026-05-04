package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/yushafro/effective-mobile-tz/internal/subscription"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
)

// Config represents the root application configuration structure.
// It combines configuration for all subsystems including subscription logic,
// PostgreSQL database, and logging.
type Config struct {
	Subscription subscription.Config `envPrefix:"APP_"`
	Postgres     postgres.Config     `envPrefix:"POSTGRES_"`
	Logger       logger.Config       `envPrefix:"LOGGER_"`
}

// LoadConfig parses environment variables into Config using env tags.
// It returns an error if parsing fails or required variables are missing.
func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
