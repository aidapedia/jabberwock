package authenticated

import (
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
)

type TokenResponse struct {
	ID           string
	Type         string
	AccessToken  string
	RefreshToken string
}

type CheckAccessTokenPayload struct {
	Token     string
	ElementID string
}

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`

	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

type LoginResponse struct {
	TokenType    string        `json:"token_type"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         userRepo.User `json:"user"`
}

func (e *LoginResponse) Transform(token TokenResponse, user userRepo.User) {
	e.TokenType = token.Type
	e.AccessToken = token.AccessToken
	e.RefreshToken = token.RefreshToken
	e.User = user
}

type LogoutRequest struct {
	Token string
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
