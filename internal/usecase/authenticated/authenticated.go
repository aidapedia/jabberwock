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
	pkgJWT "github.com/aidapedia/jabberwock/pkg/jwt"
	pkgLog "github.com/aidapedia/jabberwock/pkg/log"
)

// CheckAccessToken checks if the access token is valid
func (uc *Usecase) CheckAccessToken(ctx context.Context, payload CheckAccessTokenPayload) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthUsecase/CheckAccessToken")
	defer span.Finish(err)

	claims, err := pkgJWT.VerifyToken(strings.TrimPrefix(string(payload.Token), "Bearer "))
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

	method, path := ParseElementID(payload.ElementID)
	result, err := uc.enforcer.Enforce(role, method, path)
	if err != nil {
		return gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	if !result {
		return gers.NewWithMetadata(fmt.Errorf("user %d is not authorized to access %s", user.ID, payload.ElementID), pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	return nil
}
