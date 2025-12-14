package model

import (
	"errors"
	"net/http"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/util"
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	policyUC "github.com/aidapedia/jabberwock/internal/usecase/policy"
	"github.com/gofiber/fiber/v3"
)

type AddResourceRequest struct {
	Type   string `json:"type" validate:"oneof=http grpc"`
	Method string `json:"method"`
	Path   string `json:"path" validate:"required"`
}

func (e *AddResourceRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.AddResourceRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return policyUC.AddResourceRequest{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.Type = policyRepo.ServiceType(e.Type)
	resp.Method = e.Method
	resp.Path = e.Path
	return resp, nil
}

type AddPermissionRequest struct {
	Name              string  `json:"name" validate:"required"`
	Description       string  `json:"description" validate:"required"`
	AssignToResources []int64 `json:"assign_to_resources" `
}

func (e *AddPermissionRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.AddPermissionRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return policyUC.AddPermissionRequest{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.Name = e.Name
	resp.Description = e.Description
	resp.AssignToResources = e.AssignToResources
	return resp, nil
}

type AddRoleRequest struct {
	Name                string  `json:"name" validate:"required"`
	Description         string  `json:"description" validate:"required"`
	AssignToPermissions []int64 `json:"assign_to_permissions"`
}

func (e *AddRoleRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.AddRoleRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return policyUC.AddRoleRequest{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.Name = e.Name
	resp.Description = e.Description
	resp.AssignToPermissions = e.AssignToPermissions
	return resp, nil
}

type UpdateResourceRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Type   string `json:"type" validate:"oneof=http grpc"`
	Method string `json:"method"`
	Path   string `json:"path" validate:"required"`
}

func (e *UpdateResourceRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.UpdateResourceRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return policyUC.UpdateResourceRequest{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.ID = e.ID
	resp.Type = policyRepo.ServiceType(e.Type)
	resp.Method = e.Method
	resp.Path = e.Path
	return resp, nil
}

func (e *UpdatePermissionRequest) ToUsecase(c fiber.Ctx) policyUC.UpdatePermissionRequest {
	return policyUC.UpdatePermissionRequest{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}

type UpdatePermissionRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (e *UpdatePermissionRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.UpdatePermissionRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return policyUC.UpdatePermissionRequest{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.ID = e.ID
	resp.Name = e.Name
	resp.Description = e.Description
	return resp, nil
}

type UpdateRoleRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (e *UpdateRoleRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.UpdateRoleRequest, err error) {
	if err := c.Bind().Body(e); err != nil {
		return policyUC.UpdateRoleRequest{}, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.ID = e.ID
	resp.Name = e.Name
	resp.Description = e.Description
	return resp, nil
}

type GetUserPermissionsResponse []string

func (e *GetUserPermissionsResponse) ToSuccessResponse(resp policyUC.GetUserPermissionsResponse) *ghttp.SuccessResponse {
	if e == nil {
		*e = make([]string, len(resp.Permissions))
	}
	for _, v := range resp.Permissions {
		*e = append(*e, v.Name)
	}
	return &ghttp.SuccessResponse{
		Data: e,
	}
}

type DeleteResourceRequest struct {
	ID int64
}

func (e *DeleteResourceRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.DeleteResourceRequest, err error) {
	id := util.ToInt64(c.Params("id"))
	if id == 0 {
		return resp, gers.NewWithMetadata(errors.New("invalid id"), ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.ID = id
	return resp, nil
}

type DeletePermissionRequest struct {
	ID int64
}

func (e *DeletePermissionRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.DeletePermissionRequest, err error) {
	id := util.ToInt64(c.Params("id"))
	if id == 0 {
		return resp, gers.NewWithMetadata(errors.New("invalid id"), ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.ID = id
	return resp, nil
}

type DeleteRoleRequest struct {
	ID int64
}

func (e *DeleteRoleRequest) BindAndValidate(c fiber.Ctx) (resp policyUC.DeleteRoleRequest, err error) {
	id := util.ToInt64(c.Params("id"))
	if id == 0 {
		return resp, gers.NewWithMetadata(errors.New("invalid id"), ghttp.Metadata(http.StatusBadRequest, "Bad Request"))
	}
	resp.ID = id
	return resp, nil
}
