package session

import "context"

type Interface interface {
	FindByAccessToken(ctx context.Context, accessToken string, userID int64) (Session, error)
}
