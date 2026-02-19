package subscription

import "context"

type Service interface {
	create(ctx context.Context, sub SubReq) (SubResp, error)
	get(ctx context.Context, id string) (SubResp, error)
	update(ctx context.Context, id string, sub SubReq) (SubResp, error)
	delete(ctx context.Context, id string) error
	list(ctx context.Context, list ListSubReq) (ListSubResp, error)
	sum(ctx context.Context, subSum SubSumReq) (SubSumResp, error)
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
