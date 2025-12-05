package model

import (
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
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
}

func (e *LoginResponse) FromUsecase(resp authUC.LoginResponse) {
	e.AccessToken = resp.AccessToken
	e.RefreshToken = resp.RefreshToken
	e.User.ID = resp.User.ID
	e.User.Name = resp.User.Name
	e.User.ImageURL = resp.User.AvatarURL
	e.User.Phone = resp.User.Phone
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

type AddResourceRequest struct {
	Type   string `json:"type" validate:"oneof=http grpc"`
	Method string `json:"method"`
	Path   string `json:"path" validate:"required"`
}

func (e *AddResourceRequest) ToUsecase(c fiber.Ctx) authUC.AddResourceRequest {
	return authUC.AddResourceRequest{
		Type:   policyRepo.ServiceType(e.Type),
		Method: e.Method,
		Path:   e.Path,
	}
}

type AddPermissionRequest struct {
	Name              string  `json:"name" validate:"required"`
	Description       string  `json:"description" validate:"required"`
	AssignToResources []int64 `json:"assign_to_resources" `
}

func (e *AddPermissionRequest) ToUsecase(c fiber.Ctx) authUC.AddPermissionRequest {
	return authUC.AddPermissionRequest{
		Name:              e.Name,
		Description:       e.Description,
		AssignToResources: e.AssignToResources,
	}
}

type AddRoleRequest struct {
	Name                string  `json:"name" validate:"required"`
	Description         string  `json:"description" validate:"required"`
	AssignToPermissions []int64 `json:"assign_to_permissions"`
}

func (e *AddRoleRequest) ToUsecase(c fiber.Ctx) authUC.AddRoleRequest {
	return authUC.AddRoleRequest{
		Name:                e.Name,
		Description:         e.Description,
		AssignToPermissions: e.AssignToPermissions,
	}
}

type UpdateResourceRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Type   string `json:"type" validate:"oneof=http grpc"`
	Method string `json:"method"`
	Path   string `json:"path" validate:"required"`
}

func (e *UpdateResourceRequest) ToUsecase(c fiber.Ctx) authUC.UpdateResourceRequest {
	return authUC.UpdateResourceRequest{
		ID:     e.ID,
		Type:   policyRepo.ServiceType(e.Type),
		Method: e.Method,
		Path:   e.Path,
	}
}

type UpdatePermissionRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (e *UpdatePermissionRequest) ToUsecase(c fiber.Ctx) authUC.UpdatePermissionRequest {
	return authUC.UpdatePermissionRequest{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}

type UpdateRoleRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (e *UpdateRoleRequest) ToUsecase(c fiber.Ctx) authUC.UpdateRoleRequest {
	return authUC.UpdateRoleRequest{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}
