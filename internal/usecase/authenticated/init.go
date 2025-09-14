package authenticated

import (
	"context"

	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
	"github.com/casbin/casbin/v2"
)

type Interface interface {
	CheckAccessToken(ctx context.Context, payload CheckAccessTokenPayload) error
}

type Usecase struct {
	sessionRepo sessionRepo.Interface
	userRepo    userRepo.Interface
	enforcer    *casbin.Enforcer
}

func New() Interface {
	return &Usecase{}
}
