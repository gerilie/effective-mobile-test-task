package subscription

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository interface {
	create(ctx context.Context, sub SubReq) (SubResp, error)
	get(ctx context.Context, id string) (SubResp, error)
	update(ctx context.Context, id string, sub SubReq) (SubResp, error)
	delete(ctx context.Context, id string) error
	sum(ctx context.Context, subSum SubSumReq) (SubSumResp, error)
}

type pgRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewPGRepository(db *pgxpool.Pool) *pgRepository {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &pgRepository{
		db:      db,
		builder: builder,
	}
}
