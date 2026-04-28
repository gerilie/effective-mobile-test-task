package subscription

import (
	"context"
)

type Service interface {
	create(ctx context.Context, dto SubReq) (SubResp, error)
	get(ctx context.Context, id string) (SubResp, error)
	update(ctx context.Context, id string, dto UpdateSubReq) (SubResp, error)
	delete(ctx context.Context, id string) error
	list(ctx context.Context, dto SubListReq) (SubListResp, error)
	sum(ctx context.Context, dto SubSumReq) (SubSumResp, error)
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
