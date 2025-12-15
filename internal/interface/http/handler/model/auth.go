package model

import (
	"net/http"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
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

func (e *LoginResponse) ToSuccessResponse(resp authUC.LoginResponse) *ghttp.SuccessResponse {
	e.AccessToken = resp.AccessToken
	e.RefreshToken = resp.RefreshToken
	e.User.ID = resp.User.ID
	e.User.Name = resp.User.Name
	e.User.ImageURL = resp.User.AvatarURL
	e.User.Phone = resp.User.Phone
	for _, v := range resp.Permissions {
		e.Permissions = append(e.Permissions, v.Name)
	}
	return &ghttp.SuccessResponse{
		Data: e,
	}
}

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func (e *LoginRequest) BindAndValidate(c fiber.Ctx) (ucReq authUC.LoginRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return ucReq, ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}
	ucReq.Identity = e.Identity
	ucReq.Password = e.Password
	ucReq.IP = c.IP()
	ucReq.UserAgent = c.Get(fiber.HeaderUserAgent)
	return ucReq, nil
}

type RefreshTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (e *RefreshTokenResponse) ToSuccessResponse(resp authUC.RefreshTokenResponse) *ghttp.SuccessResponse {
	e.TokenType = resp.TokenType
	e.AccessToken = resp.AccessToken
	e.RefreshToken = resp.RefreshToken
	return &ghttp.SuccessResponse{
		Data: e,
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (e *RefreshTokenRequest) BindAndValidate(c fiber.Ctx) (ucReq authUC.RefreshTokenRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return ucReq, ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}
	ucReq.RefreshToken = e.RefreshToken
	ucReq.IP = c.IP()
	ucReq.UserAgent = c.Get(fiber.HeaderUserAgent)
	return ucReq, nil
}

type LogoutRequest struct {
	Token string `json:"token"`
}

func (e *LogoutRequest) BindAndValidate(c fiber.Ctx) (ucReq authUC.LogoutRequest, err error) {
	ucReq.Token = c.Get(fiber.HeaderAuthorization)
	return ucReq, nil
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (e *RegisterRequest) BindAndValidate(c fiber.Ctx) (ucReq authUC.RegisterRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return ucReq, ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}
	ucReq.Name = e.Name
	ucReq.Phone = e.Phone
	ucReq.Email = e.Email
	ucReq.Password = e.Password
	return ucReq, nil
}
