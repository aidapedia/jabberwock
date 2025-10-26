package authenticated

import (
	"context"
	"time"

	"github.com/aidapedia/jabberwock/pkg/config"
	pkgJWT "github.com/aidapedia/jabberwock/pkg/jwt"
	"github.com/google/uuid"
)

func (uc *Usecase) generateToken(ctx context.Context, userID int64, roleID string) (resp TokenResponse, err error) {
	id := uuid.New()
	cfg := config.GetConfig(ctx)
	// Generate access token
	accessToken, err := pkgJWT.SignToken(map[string]interface{}{
		"jti":  id.String(),
		"sub":  userID,
		"role": roleID,
		"iss":  cfg.App.Auth.Issuer,
		"exp":  time.Now().Add(time.Minute * 15).Unix(), // 15 minutes for access token
		"iat":  time.Now().Unix(),
	})
	if err != nil {
		return resp, err
	}
	// Generate refresh token
	refreshToken, err := pkgJWT.SignToken(map[string]interface{}{
		"jti":  id.String(),
		"sub":  userID,
		"role": roleID,
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
