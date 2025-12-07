package model

import (
	authUC "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/gofiber/fiber/v3"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
		Phone    string `json:"phone"`
	} `json:"user"`
	Permissions []string `json:"permissions"`
}

func (e *LoginResponse) FromUsecase(resp authUC.LoginResponse) {
	e.AccessToken = resp.AccessToken
	e.RefreshToken = resp.RefreshToken
	e.User.ID = resp.User.ID
	e.User.Name = resp.User.Name
	e.User.ImageURL = resp.User.AvatarURL
	e.User.Phone = resp.User.Phone
	for _, v := range resp.Permissions {
		e.Permissions = append(e.Permissions, v.Name)
	}
}

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func (e *LoginRequest) ToUsecase(c fiber.Ctx) authUC.LoginRequest {
	return authUC.LoginRequest{
		Identity:  e.Identity,
		Password:  e.Password,
		IP:        c.IP(),
		UserAgent: c.Get(fiber.HeaderUserAgent),
	}
}

type RefreshTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (e *RefreshTokenResponse) FromUsecase(resp authUC.RefreshTokenResponse) {
	e.TokenType = resp.TokenType
	e.AccessToken = resp.AccessToken
	e.RefreshToken = resp.RefreshToken
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (e *RefreshTokenRequest) ToUsecase(c fiber.Ctx) authUC.RefreshTokenRequest {
	return authUC.RefreshTokenRequest{
		RefreshToken: c.Get(fiber.HeaderAuthorization),
		IP:           c.IP(),
		UserAgent:    c.Get(fiber.HeaderUserAgent),
	}
}

type LogoutRequest struct {
	Token string `json:"token"`
}

func (e *LogoutRequest) ToUsecase(c fiber.Ctx) authUC.LogoutRequest {
	return authUC.LogoutRequest{
		Token: c.Get(fiber.HeaderAuthorization),
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (e *RegisterRequest) ToUsecase(c fiber.Ctx) authUC.RegisterRequest {
	return authUC.RegisterRequest{
		Name:     e.Name,
		Phone:    e.Phone,
		Email:    e.Email,
		Password: e.Password,
	}
}
