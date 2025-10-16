package model

import (
	authUC "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID              int64  `json:"id"`
		Name            string `json:"name"`
		ImageURL        string `json:"image_url"`
		Phone           string `json:"phone"`
		IsPhoneVerified bool   `json:"is_phone_verified"`
	} `json:"user"`
}

func (e *LoginResponse) FromUsecase(resp authUC.LoginResponse) {
	e.AccessToken = resp.AccessToken
	e.RefreshToken = resp.RefreshToken
	e.User.ID = resp.User.ID
	e.User.Name = resp.User.Name
	e.User.ImageURL = resp.User.ImageURL
	e.User.Phone = resp.User.Phone
	e.User.IsPhoneVerified = resp.User.IsPhoneVerified
}
