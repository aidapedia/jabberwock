package otp

import (
	"context"
	"fmt"

	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/bytedance/sonic"
	"github.com/kurniajigunawan/homestay/internal/common/cache"
	"github.com/kurniajigunawan/homestay/pkg/config"
)

func (r *Repository) SetRegistrationOTP(ctx context.Context, req RegistrationOTP) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "OTPRepository/SetRegistrationOTP")
	defer span.Finish(err)

	cfg := config.GetConfig(ctx)
	key := fmt.Sprintf(cache.RedisKeyRegistrationOTP, req.RegistrationData.Phone)
	bodyJSON, err := sonic.MarshalString(req)
	if err != nil {
		return err
	}
	err = r.redis.SET(ctx, key, bodyJSON, cfg.App.Registration.TTL)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRegistrationOTP(ctx context.Context, phone string) (result RegistrationOTP, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "OTPRepository/GetRegistrationOTP")
	defer span.Finish(err)

	key := fmt.Sprintf(cache.RedisKeyRegistrationOTP, phone)
	val, err := r.redis.GET(ctx, key)
	if err != nil {
		return RegistrationOTP{}, err
	}

	err = sonic.UnmarshalString(val, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *Repository) DeleteRegistrationOTP(ctx context.Context, phone string) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "OTPRepository/DeleteRegistrationOTP")
	defer span.Finish(err)

	key := fmt.Sprintf(cache.RedisKeyRegistrationOTP, phone)
	err = r.redis.DEL(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
