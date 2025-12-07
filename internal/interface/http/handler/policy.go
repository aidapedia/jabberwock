package handler

import (
	"errors"
	"net/http"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/gdk/util"
	"github.com/aidapedia/jabberwock/internal/common/constant"
	"github.com/aidapedia/jabberwock/internal/interface/http/handler/model"
	policyUC "github.com/aidapedia/jabberwock/internal/usecase/policy"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) AddResource(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddResource")
	defer span.Finish(nil)

	var (
		req model.AddResourceRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	err := h.policyUsecase.AddResource(ctx, req.ToUsecase(c))
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

	err := h.policyUsecase.AddPermission(ctx, req.ToUsecase(c))
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

	err := h.policyUsecase.AddRole(ctx, req.ToUsecase(c))
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

	err := h.policyUsecase.DeleteResource(ctx, policyUC.DeleteResourceRequest{
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

	err := h.policyUsecase.DeletePermission(ctx, policyUC.DeletePermissionRequest{
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

	err := h.policyUsecase.DeleteRole(ctx, policyUC.DeleteRoleRequest{
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

	err := h.policyUsecase.UpdateResource(ctx, req.ToUsecase(c))
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

	err := h.policyUsecase.UpdatePermission(ctx, req.ToUsecase(c))
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

	err := h.policyUsecase.UpdateRole(ctx, req.ToUsecase(c))
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) GetUserPermissions(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/GetUserPermissions")
	defer span.Finish(nil)

	var (
		userPermissions model.GetUserPermissionsResponse
	)

	id := util.ToInt64(c.Locals(constant.ContextKeyUserID))
	if id == 0 {
		err := errors.New("invalid id")
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	resp, err := h.policyUsecase.GetUserPermissions(ctx, id)
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	userPermissions.FromUsecase(resp)
	return ghttp.JSONResponse(c, &ghttp.SuccessResponse{
		Data: userPermissions,
	}, nil)
}
