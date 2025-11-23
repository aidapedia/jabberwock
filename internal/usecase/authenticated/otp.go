package authenticated

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	gconcurrency "github.com/aidapedia/gdk/concurrency"
	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/gdk/telemetry/tracer"
	cerror "github.com/kurniajigunawan/homestay/internal/common/error"
	otpRepo "github.com/kurniajigunawan/homestay/internal/repository/otp"
	policyRepo "github.com/kurniajigunawan/homestay/internal/repository/policy"
	whatsappRepo "github.com/kurniajigunawan/homestay/internal/repository/thirdparty/whatsapp"
	userRepo "github.com/kurniajigunawan/homestay/internal/repository/user"
	"github.com/kurniajigunawan/homestay/pkg/config"
	"go.uber.org/zap"
)

// ResendOTPRegistration resend otp verification
func (uc *Usecase) ResendOTPRegistration(ctx context.Context, req ResendOTPRegistrationRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/ResendOTPRegistration")
	defer span.Finish(err)

	otp, err := uc.otpRepo.GetRegistrationOTP(ctx, req.Phone)
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	cfg := config.GetConfig(ctx)
	// check is after 5 minutes than timenow
	subs := cfg.App.Registration.OTPResendInterval + time.Until(otp.RequestAt)
	if subs > 0 {
		return gers.NewWithMetadata(errors.New("otp is not expired"), ghttp.Metadata(http.StatusInternalServerError, fmt.Sprintf("We already resend the OTP to your phone number. Please wait %s seconds to resend", subs)))
	}

	switch req.Method {
	case MethodWhatsappText:
		return uc.sendOTPVerification(ctx, SendOTPVerificationRequest{
			Phone:    otp.RegistrationData.Phone,
			Name:     otp.RegistrationData.Name,
			Password: otp.RegistrationData.Password,
			Method:   req.Method,
		})

	}
	return nil
}

func (uc *Usecase) VerifyOTPRegistration(ctx context.Context, req VerifyOTPRegistrationRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/VerifyOTPRegistration")
	defer span.Finish(err)

	otp, err := uc.otpRepo.GetRegistrationOTP(ctx, req.Phone)
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if otp.OTP != req.OTP {
		return gers.NewWithMetadata(errors.New("invalid otp"), ghttp.Metadata(http.StatusBadRequest, "Invalid OTP"))
	}

	newUser := &userRepo.User{
		Name:       otp.RegistrationData.Name,
		Phone:      otp.RegistrationData.Phone,
		Password:   otp.RegistrationData.Password,
		Status:     userRepo.StatusActive,
		IsVerified: userRepo.VerifiedPhone,
	}
	err = uc.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	err = uc.policyRepo.AssignRole(ctx, newUser.ID, policyRepo.MemberRole)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	err = uc.otpRepo.DeleteRegistrationOTP(ctx, req.Phone)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

func (uc *Usecase) sendOTPVerification(ctx context.Context, req SendOTPVerificationRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/SendOTPVerification")
	defer span.Finish(err)

	otp, err := GenerateOTP(6)
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	err = uc.otpRepo.SetRegistrationOTP(ctx, otpRepo.RegistrationOTP{
		RegistrationData: otpRepo.RegistrationData{
			Phone:    req.Phone,
			Name:     req.Name,
			Password: req.Password,
		},
		OTP:       otp,
		RequestAt: time.Now(),
	})
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	cfg := config.GetConfig(ctx)
	switch req.Method {
	case MethodWhatsappText:
		gconcurrency.Call(ctx, func(ctx context.Context) {
			_, err = uc.whatsappRepo.SendMessageText(ctx, whatsappRepo.SendTextRequest{
				ChatID:  whatsappRepo.NewChatID(req.Phone),
				Text:    whatsappRepo.NewMessage(whatsappRepo.TemplateOTP, otp),
				Session: cfg.Secret.Whatsapp.Session,
			})
			if err != nil {
				log.ErrorCtx(ctx, "failed to send otp verification", zap.Error(err))
			}
		})
	}

	return nil
}
