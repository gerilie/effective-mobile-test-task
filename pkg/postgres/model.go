package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

// New creates and returns a new PostgreSQL connection pool.
//
// It builds a connection string from the provided configuration,
// initializes a pgx connection pool, and verifies connectivity using Ping.
//
// If the connection cannot be established or the ping fails,
// New returns an error wrapped with additional context.
func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx)

	connString := buildConnString(cfg)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	log.Info("connected to database")

	return pool, nil
}
