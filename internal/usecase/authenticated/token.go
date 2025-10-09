package authenticated

import (
	"time"

	pkgJWT "github.com/aidapedia/jabberwock/pkg/jwt"
)

func generateToken(userID int64, subID, roleID string) (access, refresh string, err error) {
	// Generate access token
	accessToken, err := pkgJWT.SignToken(map[string]interface{}{
		"sub":  subID,
		"jti":  userID,
		"role": roleID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 15 minutes for access token
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
	})
	if err != nil {
		return "", "", err
	}
	// Generate refresh token
	refreshToken, err := pkgJWT.SignToken(map[string]interface{}{
		"sub":  subID,
		"jti":  userID,
		"role": roleID,
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days for refresh token
	})
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
