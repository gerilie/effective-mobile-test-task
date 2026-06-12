package test

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	// Register pgx as a database/sql driver for Goose migrations.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
)

// PostgresContainer contains a PostgreSQL test container,
// connection string and a ready-to-use connection pool.
type PostgresContainer struct {
	Container *postgres.PostgresContainer
	DSN       string
	Pool      *pgxpool.Pool
}

// CreatePostgresContainer starts a PostgreSQL test container,
// applies database migrations and returns a ready-to-use connection pool.
//
// The caller is responsible for closing the pool and terminating
// the container after the test completes.
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:18.4-bookworm",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	err = runMigrations(ctx, dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}

	return &PostgresContainer{Container: container, DSN: dsn, Pool: pool}, nil
}

func runMigrations(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db for migrations: %w", err)
	}
	defer deferfunc.Close(ctx, db.Close, "close db for migrations")

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	err = goose.Up(db, "../../migrations")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
