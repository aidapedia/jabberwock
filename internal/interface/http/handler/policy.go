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
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) AddResource(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddResource")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.AddResourceRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.AddResource(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) AddPermission(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddPermission")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.AddPermissionRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.AddPermission(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) AddRole(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/AddRole")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.AddRoleRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.AddRole(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) DeleteResource(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/DeleteResource")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.DeleteResourceRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.DeleteResource(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) DeletePermission(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/DeletePermission")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.DeletePermissionRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.DeletePermission(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) DeleteRole(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/DeleteRole")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.DeleteRoleRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.DeleteRole(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) UpdateResource(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/UpdateResource")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.UpdateResourceRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.UpdateResource(ctx, ucReq)
	if err != nil {
		return ghttp.JSONResponse(c, nil, err)
	}

	return ghttp.JSONResponse(c, nil, nil)
}

func (h *Handler) UpdatePermission(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/UpdatePermission")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.UpdatePermissionRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.UpdatePermission(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) UpdateRole(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/UpdateRole")
	defer span.Finish(err)

	var (
		response *ghttp.SuccessResponse
		req      model.UpdateRoleRequest
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	ucReq, err := req.BindAndValidate(c)
	if err != nil {
		return
	}

	err = h.policyUsecase.UpdateRole(ctx, ucReq)
	if err != nil {
		return
	}

	return
}

func (h *Handler) GetUserPermissions(c fiber.Ctx) (err error) {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "AuthHandler/GetUserPermissions")
	defer span.Finish(err)

	var (
		response        *ghttp.SuccessResponse
		userPermissions model.GetUserPermissionsResponse
	)
	defer func() {
		err = ghttp.JSONResponse(c, response, err)
	}()

	id := util.ToInt64(c.Locals(constant.ContextKeyUserID))
	if id == 0 {
		err := errors.New("invalid id")
		return ghttp.JSONResponse(c, nil, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request")))
	}

	resp, err := h.policyUsecase.GetUserPermissions(ctx, id)
	if err != nil {
		return
	}

	response = userPermissions.ToSuccessResponse(resp)
	return
}
