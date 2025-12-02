package authenticated

import (
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
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
	Name     string
	Phone    string
	Email    string
	Password string
}

type RefreshTokenRequest struct {
	RefreshToken string

	IP        string
	UserAgent string
}

type RefreshTokenResponse struct {
	TokenType    string
	AccessToken  string
	RefreshToken string
}

type AddResourceRequest struct {
	Type   policyRepo.ServiceType
	Method string
	Path   string
}

type AddPermissionRequest struct {
	Name              string
	Description       string
	AssignToResources []int64
}

type AddRoleRequest struct {
	Name                string
	Description         string
	AssignToPermissions []int64
}
