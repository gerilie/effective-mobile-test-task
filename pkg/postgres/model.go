package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

var ErrConnect = errors.New("unable to connect to database")

type Config struct {
	Host     string `mapstructure:"postgres_host"`
	Port     string `mapstructure:"postgres_port"`
	User     string `mapstructure:"postgres_user"`
	Password string `mapstructure:"postgres_password"`
	DB       string `mapstructure:"postgres_db"`
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		net.JoinHostPort(cfg.Host, cfg.Port),
		cfg.DB,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, ErrConnect
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, ErrConnect
	}

	log.Info(ctx, "connected to database")

	return pool, nil
}
