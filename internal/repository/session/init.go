package session

import (
	"context"

	gredisengine "github.com/aidapedia/gdk/cache/engine"
)

type Interface interface {
	SetActiveSession(ctx context.Context, session Session) error
	FindByAccessToken(ctx context.Context, accessToken string, userID int64) (Session, error)

	GetLoginAttempt(ctx context.Context, userID int64) (LoginAttempt, error)
	SetLoginAttempt(ctx context.Context, userID int64, loginAttempt LoginAttempt) error
}

type Repository struct {
	redis gredisengine.Interface
}

func New(redis gredisengine.Interface) Interface {
	return &Repository{
		redis: redis,
	}
}
