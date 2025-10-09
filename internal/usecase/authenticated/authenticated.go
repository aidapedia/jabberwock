package authenticated

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	gers "github.com/aidapedia/gdk/error"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/jabberwock/internal/common"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
	pkgJWT "github.com/aidapedia/jabberwock/pkg/jwt"
	pkgLog "github.com/aidapedia/jabberwock/pkg/log"

	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"

	"github.com/google/uuid"
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
	userID, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["jti"]), 10, 64)
	if err != nil {
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	role, ok := claims["role"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("role is empty"), pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	session, err := uc.sessionRepo.FindByAccessToken(ctx, tokenID, userID)
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

	user, err := uc.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
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
	subUUID := uuid.New()
	accessToken, refreshToken, err := generateToken(user.ID, subUUID.String(), user.Type.String())
	if err != nil {
		gers.NewWithMetadata(err,
			pkgLog.Metadata(http.StatusInternalServerError, common.ErrorMessageTryAgain))
	}

	// Save session to database
	err = uc.sessionRepo.SetActiveSession(ctx, sessionRepo.Session{
		UserID:    user.ID,
		TokenID:   subUUID.String(),
		UserAgent: req.UserAgent,
		IP:        req.IP,
	})
	if err != nil {
		return LoginResponse{}, gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusInternalServerError, common.ErrorMessageTryAgain))
	}

	return LoginResponse{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: userRepo.User{
			ID:              user.ID,
			Name:            user.Name,
			ImageURL:        "https://i.imghippo.com/files/wSZE2700kcI.webp", // TODO: improve this with real avatar
			Phone:           user.Phone,
			IsPhoneVerified: user.IsPhoneVerified,
		},
	}, nil
}
