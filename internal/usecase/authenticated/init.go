package authenticated

import (
	"context"

	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"

	policyUC "github.com/aidapedia/jabberwock/internal/usecase/policy"
	"github.com/casbin/casbin/v2"
)

type Interface interface {
	// Authentication
	CheckAccessToken(ctx context.Context, req CheckAccessTokenPayload) (resp CheckAccessTokenResponse, err error)
	Login(ctx context.Context, req LoginRequest) (resp LoginResponse, err error)
	Logout(ctx context.Context, req LogoutRequest) error
	Register(ctx context.Context, req RegisterRequest) error
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (resp RefreshTokenResponse, err error)
}

type Usecase struct {
	sessionRepo sessionRepo.Interface
	userRepo    userRepo.Interface
	policyRepo  policyRepo.Interface

	policyUsecase policyUC.Interface

	enforcer *casbin.Enforcer
}

func New(policyRepo policyRepo.Interface, policyUsecase policyUC.Interface, sessionRepo sessionRepo.Interface, userRepo userRepo.Interface, enforcer *casbin.Enforcer) Interface {
	return &Usecase{
		sessionRepo:   sessionRepo,
		userRepo:      userRepo,
		enforcer:      enforcer,
		policyRepo:    policyRepo,
		policyUsecase: policyUsecase,
	}
}
