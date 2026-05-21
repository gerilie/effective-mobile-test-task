package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolConnector defines the interface for PostgreSQL pool operations.
//
// It abstracts the creation and health checking of connection pools,
// enabling dependency injection for testing and alternative implementations.
type PoolConnector interface {
	// Connect creates a new connection pool from the given connection string.
	//
	// It returns an established pool or an error if the connection parameters
	// are invalid or the database is unreachable.
	Connect(ctx context.Context, connString string) (*pgxpool.Pool, error)

	// Ping verifies that the connection pool is alive and responsive.
	//
	// It sends a lightweight check to the database and returns an error
	// if the database cannot be reached within the context deadline.
	Ping(ctx context.Context, pool *pgxpool.Pool) error
}

// poolConnector is the production implementation of PoolConnector.
//
// It delegates directly to pgxpool functions for actual database operations
// without any additional logic or middleware. The type is unexported to
// enforce creation through NewPoolConnector and prevent direct instantiation.
type poolConnector struct{}

// NewPoolConnector creates and returns a new PoolConnector instance.
//
// It returns the default production connector that uses standard pgxpool
// functions for database operations. This is the preferred way to obtain
// a PoolConnector implementation outside of tests.
func NewPoolConnector() PoolConnector {
	return &poolConnector{}
}

// Connect creates a new pgxpool.Pool using the provided connection string.
//
// It delegates directly to pgxpool.New, passing the context and connection
// string without modification. If the database is unreachable or the
// connection parameters are invalid, it returns an error from pgxpool.
func (c *poolConnector) Connect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, connString)
}

// Ping verifies connectivity by delegating to the pool's Ping method.
//
// It sends a lightweight health check to the database and returns an error
// if the connection is dead or the database does not respond within the
// context deadline. This method is a thin wrapper around pgxpool.Pool.Ping.
func (c *poolConnector) Ping(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.Ping(ctx)
}
