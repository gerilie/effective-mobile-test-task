package main

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yushafro/effective-mobile-tz/internal/config"
	"github.com/yushafro/effective-mobile-tz/internal/subscription"
	"github.com/yushafro/effective-mobile-tz/pkg/env"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
)

func initConfig(ctx context.Context) (*config.Config, error) {
	pflag.Parse()
	viper := viper.New()

	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return &config.Config{}, fmt.Errorf("bind flags: %w", err)
	}

	sub, err := readServerConfig(ctx, viper)
	if err != nil {
		return nil, fmt.Errorf("server config: %w", err)
	}

	pg, err := readPostgresConfig(ctx, viper)
	if err != nil {
		return nil, fmt.Errorf("postgres config: %w", err)
	}

	logger, err := readLoggerConfig(ctx, viper)
	if err != nil {
		return nil, fmt.Errorf("logger config: %w", err)
	}

	return &config.Config{
		Subscription: sub,
		Postgres:     pg,
		Logger:       logger,
	}, nil
}

func readServerConfig(ctx context.Context, viper *viper.Viper) (subscription.Config, error) {
	log := logger.FromContext(ctx)

	environment := viper.GetString("env")
	var configName string
	switch environment {
	case env.Host:
		configName = "host-server"
	case env.Dev:
		configName = "dev-server"
	case env.Prod:
		configName = "prod-server"
	default:
		return subscription.Config{}, fmt.Errorf("unknown env: %s", environment)
	}
	viper.SetConfigName(configName)

	viper.SetDefault("READ_HEADER_TIMEOUT", "5s")
	viper.SetDefault("READ_TIMEOUT", "15s")
	viper.SetDefault("WRITE_TIMEOUT", "30s")
	viper.SetDefault("IDLE_TIMEOUT", "120s")

	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_SECOND", "10")
	viper.SetDefault("RATE_LIMIT_BURST", "30")

	err := viper.ReadInConfig()
	if err != nil {
		return subscription.Config{}, fmt.Errorf("read file %s: %w", configName, err)
	}
	log.Info(ctx, "app config file loaded")

	cfg := subscription.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return subscription.Config{}, fmt.Errorf("decode file %s: %w", configName, err)
	}

	return cfg, nil
}

func readPostgresConfig(ctx context.Context, viper *viper.Viper) (postgres.Config, error) {
	log := logger.FromContext(ctx)

	environment := viper.GetString("env")
	var configName string
	switch environment {
	case env.Host:
		configName = "host-postgres"
	case env.Dev:
		configName = "dev-postgres"
	case env.Prod:
		configName = "prod-postgres"
	default:
		return postgres.Config{}, fmt.Errorf("unknown env: %s", environment)
	}
	viper.SetConfigName(configName)

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "5432")
	viper.SetDefault("user", "user")
	viper.SetDefault("password", "user")
	viper.SetDefault("db", "subscription")

	err := viper.ReadInConfig()
	if err != nil {
		return postgres.Config{}, fmt.Errorf("read file %s: %w", configName, err)
	}
	log.Info(ctx, "postgres config file loaded")

	cfg := postgres.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return postgres.Config{}, fmt.Errorf("decode file %s: %w", configName, err)
	}

	return cfg, nil
}

func readLoggerConfig(ctx context.Context, viper *viper.Viper) (logger.Config, error) {
	log := logger.FromContext(ctx)

	environment := viper.GetString("env")
	var configName string
	switch environment {
	case env.Dev, env.Host:
		configName = "dev-logger"
	case env.Prod:
		configName = "prod-logger"
	default:
		return logger.Config{}, fmt.Errorf("unknown env: %s", environment)
	}
	viper.SetConfigName(configName)

	viper.SetDefault("level", "info")

	err := viper.ReadInConfig()
	if err != nil {
		return logger.Config{}, fmt.Errorf("read file %s: %w", configName, err)
	}
	log.Info(ctx, "logger config file loaded")

	cfg := logger.Config{
		Env: environment,
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return logger.Config{}, fmt.Errorf("decode file %s: %w", configName, err)
	}

	return cfg, nil
}
