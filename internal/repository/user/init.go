package user

import "context"

type Interface interface {
	FindByID(ctx context.Context, id int64) (User, error)
}
