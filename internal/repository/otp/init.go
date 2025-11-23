package otp

import (
	"context"

	gredisengine "github.com/aidapedia/gdk/cache/engine"
)

type Interface interface {
	SetRegistrationOTP(ctx context.Context, req RegistrationOTP) error
	GetRegistrationOTP(ctx context.Context, phone string) (RegistrationOTP, error)
	DeleteRegistrationOTP(ctx context.Context, phone string) error
}

type Repository struct {
	redis gredisengine.Interface
}

func New(redis gredisengine.Interface) Interface {
	return &Repository{
		redis: redis,
	}
}
