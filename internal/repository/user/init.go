package user

import (
	"context"
	"database/sql"
)

type Interface interface {
	FindByID(ctx context.Context, id int64) (resp User, err error)
	FindByPhone(ctx context.Context, phone string) (resp User, err error)
	FindByEmail(ctx context.Context, email string) (resp User, err error)
	UpdateStatus(ctx context.Context, id int64, status Status) error
	CreateUser(ctx context.Context, user *User) error
}

type Repository struct {
	database *sql.DB
}

func New(database *sql.DB) Interface {
	return &Repository{
		database: database,
	}
}
