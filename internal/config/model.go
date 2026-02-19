package config

import (
	"errors"

	"github.com/yushafro/effective-mobile-tz/internal/subscription"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

type Config struct {
	Subscription subscription.Config
	Postgres     postgres.Config
	Logger       logger.Config
}

var ErrUnknownEnv = errors.New("unknown env")
