package handler

import (
	"github.com/aidapedia/jabberwock/internal/interface/http/handler/model"
	"github.com/gofiber/fiber/v3"

	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
)

func (h *Handler) Login(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Login")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.LoginRequest
		resp     model.LoginResponse
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	ucResp, err := h.authUsecase.Login(ctx, ucReq)
	if err != nil {
		return
	}

	response = resp.ToSuccessResponse(ucResp)
	return nil
}

func (h *Handler) Logout(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Logout")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.LogoutRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.authUsecase.Logout(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) Register(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Register")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.RegisterRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.authUsecase.Register(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) RefreshToken(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/RefreshToken")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.RefreshTokenRequest
		resp     model.RefreshTokenResponse
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	ucResp, err := h.authUsecase.RefreshToken(ctx, ucReq)
	if err != nil {
		return
	}

	response = resp.ToSuccessResponse(ucResp)
	return
}
