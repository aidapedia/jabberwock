package authenticated

import (
	"context"

	"github.com/casbin/casbin/v2"
	otpRepo "github.com/kurniajigunawan/homestay/internal/repository/otp"
	policyRepo "github.com/kurniajigunawan/homestay/internal/repository/policy"
	sessionRepo "github.com/kurniajigunawan/homestay/internal/repository/session"
	whatsappRepo "github.com/kurniajigunawan/homestay/internal/repository/thirdparty/whatsapp"
	userRepo "github.com/kurniajigunawan/homestay/internal/repository/user"

	policyUsecase "github.com/kurniajigunawan/homestay/internal/usecase/policy"
)

type Interface interface {
	// Authentication
	CheckAccessToken(ctx context.Context, req CheckAccessTokenPayload) (resp CheckAccessTokenResponse, err error)
	Login(ctx context.Context, req LoginRequest) (resp LoginResponse, err error)
	Logout(ctx context.Context, req LogoutRequest) error
	Register(ctx context.Context, req RegisterRequest) error
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (resp RefreshTokenResponse, err error)

	ResendOTPRegistration(ctx context.Context, req ResendOTPRegistrationRequest) (err error)
	VerifyOTPRegistration(ctx context.Context, req VerifyOTPRegistrationRequest) (err error)
}

type Usecase struct {
	whatsappRepo whatsappRepo.Interface
	sessionRepo  sessionRepo.Interface
	otpRepo      otpRepo.Interface
	policyRepo   policyRepo.Interface
	userRepo     userRepo.Interface

	policyUsecase policyUsecase.Interface
	enforcer      *casbin.Enforcer
}

func New(whatsappRepo whatsappRepo.Interface, policyRepo policyRepo.Interface, sessionRepo sessionRepo.Interface, userRepo userRepo.Interface,
	otpRepo otpRepo.Interface, policyUsecase policyUsecase.Interface, enforcer *casbin.Enforcer) Interface {
	return &Usecase{
		whatsappRepo: whatsappRepo,
		policyRepo:   policyRepo,
		sessionRepo:  sessionRepo,
		otpRepo:      otpRepo,
		userRepo:     userRepo,
		enforcer:     enforcer,

		policyUsecase: policyUsecase,
	}
}
