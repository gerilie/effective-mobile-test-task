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

	log := logger.FromContext(ctx)
	viper := viper.New()

	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Error(ctx, "binding PFlags", zap.Error(err))

		return &config.Config{}, err
	}

	sub, err := readServerConfig(ctx, viper)
	if err != nil {
		return nil, err
	}

	pg, err := readPostgresConfig(ctx, viper)
	if err != nil {
		return nil, err
	}

	logger, err := readLoggerConfig(ctx, viper)
	if err != nil {
		return nil, err
	}

	return &config.Config{
		Subscription: sub,
		Postgres:     pg,
		Logger:       logger,
	}, nil
}

func readServerConfig(ctx context.Context, viper *viper.Viper) (subscription.Config, error) {
	log := logger.FromContext(ctx)

	switch viper.GetString("env") {
	case config.Dev:
		viper.SetConfigName("dev-server")
	case config.Prod:
		viper.SetConfigName("prod-server")
	default:
		log.Error(ctx, "unknown env", zap.String("env", viper.GetString("env")))

		return subscription.Config{}, config.ErrUnknownEnv
	}

	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("CLIENT_TIMEOUT", "2s")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(ctx, "unable to read app config", zap.Error(err))

		return subscription.Config{}, err
	}
	log.Info(ctx, "app config file loaded")

	cfg := subscription.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Error(ctx, "unable to decode app config", zap.Error(err))

		return subscription.Config{}, err
	}

	return cfg, nil
}

func readPostgresConfig(ctx context.Context, viper *viper.Viper) (postgres.Config, error) {
	log := logger.FromContext(ctx)

	switch viper.GetString("env") {
	case config.Dev:
		viper.SetConfigName("dev-postgres")
	case config.Prod:
		viper.SetConfigName("prod-postgres")
	default:
		log.Error(ctx, config.ErrUnknownEnv.Error(), zap.String("env", viper.GetString("env")))

		return postgres.Config{}, config.ErrUnknownEnv
	}

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "5432")
	viper.SetDefault("user", "user")
	viper.SetDefault("password", "user")
	viper.SetDefault("db", "subscription")

	err := viper.ReadInConfig()
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

func readLoggerConfig(ctx context.Context, viper *viper.Viper) (logger.Config, error) {
	log := logger.FromContext(ctx)

	switch viper.GetString("env") {
	case config.Dev:
		viper.SetConfigName("dev-logger")
	case config.Prod:
		viper.SetConfigName("prod-logger")
	default:
		log.Error(ctx, config.ErrUnknownEnv.Error(), zap.String("env", viper.GetString("env")))

		return logger.Config{}, config.ErrUnknownEnv
	}

	viper.SetDefault("level", "info")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(ctx, "unable to read logger config", zap.Error(err))

		return logger.Config{}, err
	}
	log.Info(ctx, "logger config file loaded")

	cfg := logger.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Error(ctx, "unable to decode logger config", zap.Error(err))

		return logger.Config{}, err
	}

	return cfg, nil
}
