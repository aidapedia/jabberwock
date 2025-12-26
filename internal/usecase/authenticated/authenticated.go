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
	"github.com/aidapedia/gdk/util"
	gvalidation "github.com/aidapedia/gdk/validation"
	cerror "github.com/aidapedia/jabberwock/internal/common/error"
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"

	policyUC "github.com/aidapedia/jabberwock/internal/usecase/policy"
)

// CheckAccessToken checks if the access token is valid
func (uc *Usecase) CheckAccessToken(ctx context.Context, req CheckAccessTokenPayload) (resp CheckAccessTokenResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/CheckAccessToken")
	defer span.Finish(err)

	claims, err := gjwt.VerifyToken(strings.TrimPrefix(string(req.Token), "Bearer "))
	if err != nil {
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	tokenID, ok := claims["jti"].(string)
	if !ok {
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(errors.New("jti is empty"), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}
	role, ok := claims["role"].(string)
	if !ok {
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(errors.New("role is empty"), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	session, err := uc.sessionRepo.FindActiveSessionByTokenID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, cerror.ErrorNotFound) {
			return CheckAccessTokenResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusForbidden, "Session not found"))
		}
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, "Internal Server Error"))
	}

	user, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, cerror.ErrorNotFound) {
			return CheckAccessTokenResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusForbidden, "User not found"))
		}
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if err = uc.validateUser(user); err != nil {
		return CheckAccessTokenResponse{}, err
	}

	method, path := ParseElementID(req.ElementID)
	result, err := uc.enforcer.Enforce(policyUC.StdPolicy(policyRepo.Policy{
		Role:   role,
		Path:   path,
		Type:   "http",
		Method: method,
	})...)
	if err != nil {
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	if !result {
		return CheckAccessTokenResponse{}, gers.NewWithMetadata(fmt.Errorf("user %d is not authorized to access %s", user.ID, req.ElementID), ghttp.Metadata(http.StatusUnauthorized, "Unauthorized"))
	}

	return CheckAccessTokenResponse{
		UserID: user.ID,
	}, nil
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

	permissions, err := uc.policyUsecase.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return LoginResponse{}, err
	}

	roles, err := uc.policyRepo.GetRoleByUserID(ctx, user.ID)
	if err != nil {
		return LoginResponse{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if len(roles) == 0 {
		return LoginResponse{}, gers.NewWithMetadata(errors.New("user has no role"), ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	var role string
	for _, r := range roles {
		if role != "" {
			role += ","
		}
		role += r.Name
	}

	// Generate token
	tokenResp, err := uc.generateToken(ctx, user.ID, role)
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

	resp.Transform(tokenResp, user, permissions.Permissions)
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
	newUser := &userRepo.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: ghash.Hash(req.Password),
		Status:   userRepo.StatusActive,
		Email:    req.Email,
	}
	err = uc.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	// craete new user role
	err = uc.policyRepo.AssignRole(ctx, newUser.ID, policyRepo.MemberRole)
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

	token, ok := claims["jti"]
	if !ok {
		return RefreshTokenResponse{}, gers.NewWithMetadata(errors.New("invalid token"),
			ghttp.Metadata(http.StatusInternalServerError, "Your refresh token is invalid"))
	}
	tokenID := util.ToStr(token)

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
	tokenResp, err := uc.refreshToken(ctx, sessions, util.ToStr(claims["role"]))
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
