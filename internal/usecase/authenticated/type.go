package authenticated

import (
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
)

type TokenResponse struct {
	ID           string
	AccessToken  string
	RefreshToken string
}

type CheckAccessTokenPayload struct {
	Token     string
	ElementID string
}

type LoginRequest struct {
	Identity string
	Password string

	IP        string
	UserAgent string
}

type LoginResponse struct {
	TokenType    string
	AccessToken  string
	RefreshToken string
	User         userRepo.User
}

type LogoutRequest struct {
	Token string
}
