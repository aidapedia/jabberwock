package user

import "context"

type Interface interface {
	FindByID(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateStatus(ctx context.Context, id int64, status Status) error
}

type Repository struct {
}

func New() Interface {
	return &Repository{}
}
