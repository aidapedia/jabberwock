package userdatacenter

import (
	"context"

	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
)

type Interface interface {
	GetUserByID(ctx context.Context, id int64) (User, error)
}

type Usecase struct {
	userRepo userRepo.Interface
}

func New(userRepo userRepo.Interface) Interface {
	return &Usecase{
		userRepo: userRepo,
	}
}
