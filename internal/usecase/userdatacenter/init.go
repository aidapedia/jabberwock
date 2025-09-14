package userdatacenter

import "context"

type Interface interface {
	GetUserByID(ctx context.Context, id int64) (User, error)
}

type Usecase struct {
}

func New() Interface {
	return &Usecase{}
}
