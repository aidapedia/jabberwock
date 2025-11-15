package authenticated

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	ghash "github.com/aidapedia/gdk/cryptography/hash"
	gjwt "github.com/aidapedia/gdk/cryptography/jwt"
	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
	gvalidation "github.com/aidapedia/gdk/validation"
	cerror "github.com/aidapedia/jabberwock/internal/common/error"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"

	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
)

// CheckAccessToken checks if the access token is valid
func (uc *Usecase) CheckAccessToken(ctx context.Context, req CheckAccessTokenPayload) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/CheckAccessToken")
	defer span.Finish(err)

	claims, err := gjwt.VerifyToken(strings.TrimPrefix(string(req.Token), "Bearer "))
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	tokenID, ok := claims["jti"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("jti is empty"), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	role, ok := claims["role"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("role is empty"), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	session, err := uc.sessionRepo.FindActiveSessionByTokenID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, cerror.ErrorNotFound) {
			return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusForbidden, "Session not found"))
		}
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, "Internal Server Error"))
	}

	user, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, cerror.ErrorNotFound) {
			return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusForbidden, "User not found"))
		}
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if err = uc.validateUser(user); err != nil {
		return err
	}

	method, path := ParseElementID(req.ElementID)
	result, err := uc.enforcer.Enforce(role, path, method)
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	if !result {
		return gers.NewWithMetadata(fmt.Errorf("user %d is not authorized to access %s", user.ID, req.ElementID), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	return nil
}

// Login do login for user get sessions
func (uc *Usecase) Login(ctx context.Context, req LoginRequest) (resp LoginResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/Login")
	defer span.Finish(err)

	// Check if user identity is email or phone number
	var user userRepo.User
	if gvalidation.IsEmail(req.Identity) {
		user, err = uc.userRepo.FindByEmail(ctx, req.Identity)
	} else {
		user, err = uc.userRepo.FindByPhone(ctx, req.Identity)
	}
	if err != nil {
		if errors.Is(err, cerror.ErrorNotFound) {
			return LoginResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, "Cannot find your account."))
		}
		return LoginResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
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
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	// Save session to database
	err = uc.sessionRepo.CreateActiveSession(ctx, &sessionRepo.Session{
		UserID:    user.ID,
		Token:     tokenResp.ID,
		UserAgent: req.UserAgent,
		IP:        req.IP,
	})
	if err != nil {
		return LoginResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	resp.Transform(tokenResp, user)
	return resp, nil
}

func (uc *Usecase) Logout(ctx context.Context, req LogoutRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/Logout")
	defer span.Finish(err)

	claims, err := gjwt.VerifyToken(strings.TrimPrefix(req.Token, "Bearer "))
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		return gers.NewWithMetadata(errors.New("jti is empty"), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	err = uc.sessionRepo.DeleteActiveSession(ctx, tokenID)
	if err != nil {
		return gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

// Register is a function to handle register
func (uc *Usecase) Register(ctx context.Context, req RegisterRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/Register")
	defer span.Finish(err)

	// Check if user is already exist
	identities := map[string]string{
		"email": req.Email,
		"phone": req.Phone,
	}
	for key, identity := range identities {
		var exist bool
		exist, err = uc.isExistUser(ctx, identity)
		if err != nil {
			return gers.NewWithMetadata(err,
				ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
		}
		if exist {
			return gers.NewWithMetadata(errors.New("user already exist"),
				ghttp.Metadata(http.StatusBadRequest, fmt.Sprintf("%s is already registered", key)))
		}
	}
	// Create user
	err = uc.userRepo.CreateUser(ctx, &userRepo.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: ghash.Hash(req.Password),
		Type:     userRepo.TypeUser,
		Status:   userRepo.StatusActive,
		Email:    req.Email,
	})
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

// RefreshToken is a function to handle refresh token
func (uc *Usecase) RefreshToken(ctx context.Context, req RefreshTokenRequest) (resp RefreshTokenResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/RefreshToken")
	defer span.Finish(err)

	// Verify refresh token
	claims, err := gjwt.VerifyToken(req.RefreshToken)
	if err != nil {
		return RefreshTokenResponse{}, gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, "Your refresh token is invalid"))
	}

	tokenID, ok := claims["sub"].(string)
	if !ok {
		return RefreshTokenResponse{}, gers.NewWithMetadata(errors.New("invalid token"),
			ghttp.Metadata(http.StatusInternalServerError, "Your refresh token is invalid"))
	}

	// Get session from database
	sessions, err := uc.sessionRepo.FindActiveSessionByTokenID(ctx, tokenID)
	if err != nil {
		return RefreshTokenResponse{}, gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	// Validate Refresh Token
	if sessions.IP != req.IP {
		return RefreshTokenResponse{}, gers.NewWithMetadata(errors.New("ip not match"),
			ghttp.Metadata(http.StatusInternalServerError, "Your refresh token is invalid"))
	}
	if sessions.UserAgent != req.UserAgent {
		return RefreshTokenResponse{}, gers.NewWithMetadata(errors.New("user agent not match"),
			ghttp.Metadata(http.StatusInternalServerError, "Your refresh token is invalid"))
	}

	// Generate token
	tokenResp, err := uc.refreshToken(ctx, sessions, claims["role"].(string))
	if err != nil {
		return RefreshTokenResponse{}, gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	// Update session token last updated
	err = uc.sessionRepo.UpdateRefreshDateByTokenID(ctx, tokenID)
	if err != nil {
		return RefreshTokenResponse{}, gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return RefreshTokenResponse{
		TokenType:    "Bearer",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}, nil
}

func (uc *Usecase) isExistUser(ctx context.Context, identity string) (bool, error) {
	var (
		user userRepo.User
		err  error
	)
	if gvalidation.IsEmail(identity) {
		user, err = uc.userRepo.FindByEmail(ctx, identity)
	} else {
		user, err = uc.userRepo.FindByPhone(ctx, identity)
	}
	if err != nil {
		if errors.Is(err, cerror.ErrorNotFound) {
			return false, nil
		}
		return false, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}
	return !user.IsEmpty(), nil
}
