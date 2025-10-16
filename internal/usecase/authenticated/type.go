package authenticated

import (
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
)

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
