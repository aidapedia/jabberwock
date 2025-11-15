package authenticated

import (
	"context"
	"time"

	gjwt "github.com/aidapedia/gdk/cryptography/jwt"
	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	"github.com/aidapedia/jabberwock/pkg/config"
	"github.com/google/uuid"
)

func (uc *Usecase) generateToken(ctx context.Context, userID int64, roleStr string) (resp TokenResponse, err error) {
	return createToken(ctx, userID, roleStr)
}

func (uc *Usecase) refreshToken(ctx context.Context, session sessionRepo.Session, roleStr string) (resp TokenResponse, err error) {
	return createToken(ctx, session.UserID, roleStr)
}

func createToken(ctx context.Context, userID int64, roleStr string) (resp TokenResponse, err error) {
	id := uuid.New()
	cfg := config.GetConfig(ctx)
	// Generate access token
	accessToken, err := gjwt.SignToken(map[string]interface{}{
		"jti":  id.String(),
		"sub":  userID,
		"role": roleStr,
		"iss":  cfg.App.Auth.Issuer,
		"exp":  time.Now().Add(time.Minute * 15).Unix(), // 15 minutes for access token
		"iat":  time.Now().Unix(),
	})
	if err != nil {
		return resp, err
	}
	// Generate refresh token
	refreshToken, err := gjwt.SignToken(map[string]interface{}{
		"jti":  id.String(),
		"sub":  userID,
		"role": roleStr,
		"iss":  cfg.App.Auth.Issuer,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days for refresh token
		"iat":  time.Now().Unix(),
	})
	if err != nil {
		return resp, err
	}

	resp.ID = id.String()
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	return resp, nil
}
