package authenticated

import (
	"context"

	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
	"github.com/casbin/casbin/v2"
)

type Interface interface {
	// Authentication
	CheckAccessToken(ctx context.Context, req CheckAccessTokenPayload) error
	Login(ctx context.Context, req LoginRequest) (resp LoginResponse, err error)
	Logout(ctx context.Context, req LogoutRequest) error
	Register(ctx context.Context, req RegisterRequest) error
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (resp RefreshTokenResponse, err error)

	// Policy Management
	LoadPolicy(ctx context.Context, serviceType policyRepo.ServiceType) (err error)
	AddResource(ctx context.Context, req AddResourceRequest) (err error)
	AddPermission(ctx context.Context, req AddPermissionRequest) (err error)
	AddRole(ctx context.Context, req AddRoleRequest) (err error)
}

type Usecase struct {
	sessionRepo sessionRepo.Interface
	userRepo    userRepo.Interface
	policyRepo  policyRepo.Interface

	enforcer *casbin.Enforcer
}

func New(policyRepo policyRepo.Interface, sessionRepo sessionRepo.Interface, userRepo userRepo.Interface, enforcer *casbin.Enforcer) Interface {
	return &Usecase{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		enforcer:    enforcer,
		policyRepo:  policyRepo,
	}
}
