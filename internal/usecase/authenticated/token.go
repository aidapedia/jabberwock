package authenticated

import (
	"context"
	"errors"
	"time"

	pkgJWT "github.com/aidapedia/jabberwock/pkg/jwt"
	"github.com/google/uuid"
)

func (uc *Usecase) generateToken(ctx context.Context, userID int64, roleID string) (resp TokenResponse, err error) {
	id := uuid.New()
	loop := 0
	for {
		// Check if token ID is already used
		// This sanitizes the token ID to prevent any potential issues
		_, err = uc.sessionRepo.FindActiveSessionByTokenID(ctx, id.String())
		if err == nil {
			break
		}
		id = uuid.New()
		// Max 5 attempts to generate unique token ID
		loop++
		if loop > 5 {
			return resp, errors.New("failed to generate token")
		}
	}
	// Generate access token
	accessToken, err := pkgJWT.SignToken(map[string]interface{}{
		"sub":  id.String(),
		"jti":  userID,
		"role": roleID,
		"exp":  time.Now().Add(time.Minute * 15).Unix(), // 15 minutes for access token
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
	})
	if err != nil {
		return resp, err
	}
	// Generate refresh token
	refreshToken, err := pkgJWT.SignToken(map[string]interface{}{
		"sub":  id.String(),
		"jti":  userID,
		"role": roleID,
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days for refresh token
	})
	if err != nil {
		return resp, err
	}

	resp.ID = id.String()
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	return resp, nil
}
