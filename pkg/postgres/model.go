package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

// ErrConnectionFailed is returned by New when the connection to PostgreSQL
// cannot be established or when the initial ping to verify connectivity fails.
var ErrConnectionFailed = errors.New("connection failed")

// New creates and returns a new PostgreSQL connection pool.
//
// It builds a connection string from the provided configuration, then uses
// the given PoolConnector to establish a connection and verify connectivity
// with a ping. The connector parameter enables dependency injection, allowing
// mock implementations to be substituted during testing.
//
// If the connection cannot be established or the ping fails, New returns an
// error wrapped with additional context describing the failure. On success,
// it logs an informational message and returns the initialized pool.
func New(ctx context.Context, cfg Config, connector PoolConnector) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx)

	connString := BuildConnString(cfg)

	pool, err := connector.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectionFailed, err)
	}

	err = connector.Ping(ctx, pool)
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf("%w: ping failed: %w", ErrConnectionFailed, err)
	}

	log.Info("connected to database")

	return pool, nil
}
