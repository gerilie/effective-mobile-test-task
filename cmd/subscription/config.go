package main

import (
	"context"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yushafro/effective-mobile-tz/internal/config"
	"github.com/yushafro/effective-mobile-tz/internal/subscription"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/zap"
)

func initConfig(ctx context.Context) (*config.Config, error) {
	pflag.String("env", config.Dev, "Environment")
	pflag.Parse()

	sub, err := readServerConfig(ctx)
	if err != nil {
		return nil, err
	}

	pg, err := readPostgresConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &config.Config{
		Subscription: sub,
		Postgres:     pg,
	}, nil
}

func readServerConfig(ctx context.Context) (subscription.Config, error) {
	log := logger.FromContext(ctx)
	viper := viper.New()

	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Error(ctx, "binding PFlags", zap.Error(err))

		return subscription.Config{}, err
	}

	switch viper.GetString("env") {
	case config.Dev:
		viper.SetConfigName("dev-server")
	case config.Prod:
		viper.SetConfigName("prod-server")
	default:
		log.Error(ctx, "unknown env", zap.String("env", viper.GetString("env")))

		return subscription.Config{}, err
	}

	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("CLIENT_TIMEOUT", "2s")

	err = viper.ReadInConfig()
	if err != nil {
		log.Error(ctx, "unable to read server config", zap.Error(err))

		return subscription.Config{}, err
	}
	log.Info(ctx, "server config file loaded")

	cfg := subscription.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Error(ctx, "unable to decode server config", zap.Error(err))

		return subscription.Config{}, err
	}

	return cfg, nil
}

func readPostgresConfig(ctx context.Context) (postgres.Config, error) {
	log := logger.FromContext(ctx)
	viper := viper.New()

	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	viper.SetEnvPrefix("postgres")

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Error(ctx, "binding PFlags", zap.Error(err))

		return postgres.Config{}, err
	}

	switch viper.GetString("env") {
	case config.Dev:
		viper.SetConfigName("dev-postgres")
	case config.Prod:
		viper.SetConfigName("prod-postgres")
	default:
		log.Error(ctx, config.ErrUnknownEnv.Error(), zap.String("env", viper.GetString("env")))

		return postgres.Config{}, err
	}

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "5432")
	viper.SetDefault("user", "user")
	viper.SetDefault("password", "user")
	viper.SetDefault("db", "subscription")

	err = viper.ReadInConfig()
	if err != nil {
		log.Error(ctx, "unable to read postgres config", zap.Error(err))

		return postgres.Config{}, err
	}
	log.Info(ctx, "postgres config file loaded")

	cfg := postgres.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Error(ctx, "unable to decode postgres config", zap.Error(err))

		return postgres.Config{}, err
	}

	return cfg, nil
}
