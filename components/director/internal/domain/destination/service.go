package destination

import "context"

type Repository interface {
	Upsert(ctx context.Context) error
	Delete(ctx context.Context) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}
