package session

import (
	"context"
	"database/sql"

	gredisengine "github.com/aidapedia/gdk/cache/engine"
)

type Interface interface {
	CreateActiveSession(ctx context.Context, session *Session) error
	FindActiveSessionByTokenID(ctx context.Context, tokenID string) (Session, error)
	DeleteActiveSession(ctx context.Context, tokenID string) error

	GetLoginAttempt(ctx context.Context, userID int64) (LoginAttempt, error)
	SetLoginAttempt(ctx context.Context, userID int64, loginAttempt LoginAttempt) error
}

type Repository struct {
	database *sql.DB
	redis    gredisengine.Interface
}

func New(database *sql.DB, redis gredisengine.Interface) Interface {
	return &Repository{
		database: database,
		redis:    redis,
	}
}
