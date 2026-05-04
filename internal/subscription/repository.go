package subscription

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines persistence operations for subscriptions.
type Repository interface {
	create(ctx context.Context, model sub) (sub, error)
	get(ctx context.Context, id string) (sub, error)
	update(ctx context.Context, id string, model updateSub) (sub, error)
	delete(ctx context.Context, id string) error
	list(ctx context.Context, model subList) ([]sub, error)
	sum(ctx context.Context, model subSum) (subSum, error)
}

type pgRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

// NewPGRepository creates PostgreSQL implementation of Repository.
func NewPGRepository(db *pgxpool.Pool) *pgRepository {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &pgRepository{
		db:      db,
		builder: builder,
	}
}
