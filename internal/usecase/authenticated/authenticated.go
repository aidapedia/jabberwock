package authenticated

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	gers "github.com/aidapedia/gdk/error"
	"github.com/aidapedia/gdk/telemetry/tracer"
	gvalidation "github.com/aidapedia/gdk/validation"
	"github.com/aidapedia/jabberwock/internal/common"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
	pkgJWT "github.com/aidapedia/jabberwock/pkg/jwt"
	pkgLog "github.com/aidapedia/jabberwock/pkg/log"

	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
)

// CheckAccessToken checks if the access token is valid
func (uc *Usecase) CheckAccessToken(ctx context.Context, req CheckAccessTokenPayload) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthUsecase/CheckAccessToken")
	defer span.Finish(err)

	claims, err := pkgJWT.VerifyToken(strings.TrimPrefix(string(req.Token), "Bearer "))
	if err != nil {
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	tokenID, ok := claims["sub"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("sub is empty"), pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	role, ok := claims["role"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("role is empty"), pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	session, err := uc.sessionRepo.FindActiveSessionByTokenID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusForbidden, "Session not found"))
		}
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, "Internal Server Error"))
	}

	user, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusForbidden, "User not found"))
		}
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, "Internal Server Error"))
	}

	if !user.IsPhoneVerified {
		return gers.NewWithMetadata(fmt.Errorf("user is not verified by phone"), pkgLog.Metadata(http.StatusForbidden, "Your phone number is not verified"))
	}

	method, path := ParseElementID(req.ElementID)
	result, err := uc.enforcer.Enforce(role, method, path)
	if err != nil {
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	if !result {
		return gers.NewWithMetadata(fmt.Errorf("user %d is not authorized to access %s", user.ID, req.ElementID), pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	return nil
}

// Login do login for user get sessions
func (uc *Usecase) Login(ctx context.Context, req LoginRequest) (resp LoginResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthUsecase/Login")
	defer span.Finish(err)

	// Check if user identity is email or phone number
	var user userRepo.User
	if gvalidation.IsEmail(req.Identity) {
		user, err = uc.userRepo.FindByEmail(ctx, req.Identity)
	} else {
		user, err = uc.userRepo.FindByPhone(ctx, req.Identity)
	}
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return LoginResponse{}, gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, "We cannot find your account"))
		}
		return LoginResponse{}, gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, common.ErrorMessageTryAgain))
	}

	// Check if user phone number is already verified
	err = uc.validateUser(user)
	if err != nil {
		return LoginResponse{}, err
	}

	// Validate password
	err = uc.validationPassword(ctx, user, req.Password)
	if err != nil {
		return LoginResponse{}, err
	}

	// Generate token
	tokenResp, err := uc.generateToken(ctx, user.ID, user.Type.String())
	if err != nil {
		return LoginResponse{}, gers.NewWithMetadata(err,
			pkgLog.Metadata(http.StatusInternalServerError, common.ErrorMessageTryAgain))
	}

	// Save session to database
	err = uc.sessionRepo.SetActiveSession(ctx, sessionRepo.Session{
		UserID:    user.ID,
		TokenID:   tokenResp.ID,
		UserAgent: req.UserAgent,
		IP:        req.IP,
	})
	if err != nil {
		return LoginResponse{}, gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, common.ErrorMessageTryAgain))
	}

	return LoginResponse{
		TokenType:    "Bearer",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		User: userRepo.User{
			ID:              user.ID,
			Name:            user.Name,
			ImageURL:        user.ImageURL,
			Phone:           user.Phone,
			IsPhoneVerified: user.IsPhoneVerified,
		},
	}, nil
}

func (uc *Usecase) Logout(ctx context.Context, req LogoutRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthUsecase/Logout")
	defer span.Finish(err)

	claims, err := pkgJWT.VerifyToken(strings.TrimPrefix(req.Token, "Bearer "))
	if err != nil {
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	tokenID, ok := claims["sub"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("sub is empty"), pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	err = uc.sessionRepo.DeleteActiveSession(ctx, tokenID)
	if err != nil {
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, common.ErrorMessageTryAgain))
	}

	return nil
}
