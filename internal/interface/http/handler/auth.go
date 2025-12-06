package handler

import (
	"net/http"

	"github.com/aidapedia/jabberwock/internal/interface/http/handler/model"
	"github.com/gofiber/fiber/v3"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
)

func (h *Handler) Login(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Login")
	defer span.Finish(nil)

	var (
		req  model.LoginRequest
		resp model.LoginResponse
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	ucResp, err := h.authUsecase.Login(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	resp.FromUsecase(ucResp)
	return ghttp.JSONResponse(c, &ghttp.SuccessResponse{
		Data: resp,
	}, nil)
}

func (h *Handler) Logout(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Logout")
	defer span.Finish(nil)

	var (
		req model.LogoutRequest
	)
	err := h.authUsecase.Logout(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) Register(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Register")
	defer span.Finish(nil)

	var (
		req model.RegisterRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.Register(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) RefreshToken(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/RefreshToken")
	defer span.Finish(nil)

	var (
		req  model.RefreshTokenRequest
		resp model.RefreshTokenResponse
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	ucResp, err := h.authUsecase.RefreshToken(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	resp.FromUsecase(ucResp)
	return ghttp.JSONResponse(c, &ghttp.SuccessResponse{
		Data: resp,
	}, nil)
}
