package handler

import (
	"net/http"

	"github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/gofiber/fiber/v3"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"

	pkgLog "github.com/aidapedia/jabberwock/pkg/log"
)

func (h *Handler) Login(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/Login")
	defer span.Finish(nil)

	var req authenticated.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, pkgLog.Metadata(http.StatusUnauthorized, "Unauthorized")))
		return err
	}

	resp, err := h.usecase.Login(ctx, req)
	if err != nil {
		ghttp.JSONResponse(c, nil, err)
		return err
	}

	ghttp.JSONResponse(c, resp, nil)
	return nil
}
