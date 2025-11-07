package handler

import (
	"net/http"

	"github.com/aidapedia/jabberwock/internal/interface/http/handler/model"
	"github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/gofiber/fiber/v3"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
)

func (h *Handler) Login(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Login")
	defer span.Finish(nil)

	var (
		req  authenticated.LoginRequest
		resp model.LoginResponse
	)
	if err := c.Bind().Body(&req); err != nil {
		ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
		return err
	}

	ucResp, err := h.usecase.Login(ctx, req)
	if err != nil {
		ghttp.JSONResponse(c, nil, err)
		return err
	}

	resp.FromUsecase(ucResp)
	ghttp.JSONResponse(c, resp, nil)
	return nil
}

func (h *Handler) Logout(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Logout")
	defer span.Finish(nil)

	if err := h.usecase.Logout(ctx, authenticated.LogoutRequest{
		Token: string(c.Request().Header.Peek("Authorization")),
	}); err != nil {
		ghttp.JSONResponse(c, nil, err)
		return err
	}

	ghttp.JSONResponse(c, nil, nil)
	return nil
}
