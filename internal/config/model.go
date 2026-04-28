package config

import (
	"github.com/yushafro/effective-mobile-tz/internal/subscription"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
)

type Config struct {
	Subscription subscription.Config
	Postgres     postgres.Config
	Logger       logger.Config
}
