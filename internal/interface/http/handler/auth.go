package handler

import (
	"errors"
	"net/http"

	"github.com/aidapedia/jabberwock/internal/interface/http/handler/model"
	"github.com/gofiber/fiber/v3"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/gdk/util"

	authUC "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
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

func (h *Handler) AddResource(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddResource")
	defer span.Finish(nil)

	var (
		req model.AddResourceRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.AddResource(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) AddPermission(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddPermission")
	defer span.Finish(nil)

	var (
		req model.AddPermissionRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.AddPermission(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) AddRole(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddRole")
	defer span.Finish(nil)

	var (
		req model.AddRoleRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.AddRole(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) DeleteResource(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/DeleteResource")
	defer span.Finish(nil)

	id := util.ToInt64(c.Params("id"))
	if id == 0 {
		err := errors.New("invalid id")
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.DeleteResource(ctx, authUC.DeleteResourceRequest{
		ID: id,
	})
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) DeletePermission(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/DeletePermission")
	defer span.Finish(nil)

	id := util.ToInt64(c.Params("id"))
	if id == 0 {
		err := errors.New("invalid id")
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.DeletePermission(ctx, authUC.DeletePermissionRequest{
		ID: id,
	})
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) DeleteRole(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/DeleteRole")
	defer span.Finish(nil)

	id := util.ToInt64(c.Params("id"))
	if id == 0 {
		err := errors.New("invalid id")
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.DeleteRole(ctx, authUC.DeleteRoleRequest{
		ID: id,
	})
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) UpdateResource(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/UpdateResource")
	defer span.Finish(nil)

	var (
		req model.UpdateResourceRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.UpdateResource(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) UpdatePermission(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/UpdatePermission")
	defer span.Finish(nil)

	var (
		req model.UpdatePermissionRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.UpdatePermission(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) UpdateRole(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/UpdateRole")
	defer span.Finish(nil)

	var (
		req model.UpdateRoleRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.authUsecase.UpdateRole(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}
