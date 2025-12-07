package model

import (
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	policyUC "github.com/aidapedia/jabberwock/internal/usecase/policy"
	"github.com/gofiber/fiber/v3"
)

type AddResourceRequest struct {
	Type   string `json:"type" validate:"oneof=http grpc"`
	Method string `json:"method"`
	Path   string `json:"path" validate:"required"`
}

func (e *AddResourceRequest) ToUsecase(c fiber.Ctx) policyUC.AddResourceRequest {
	return policyUC.AddResourceRequest{
		Type:   policyRepo.ServiceType(e.Type),
		Method: e.Method,
		Path:   e.Path,
	}
}

type AddPermissionRequest struct {
	Name              string  `json:"name" validate:"required"`
	Description       string  `json:"description" validate:"required"`
	AssignToResources []int64 `json:"assign_to_resources" `
}

func (e *AddPermissionRequest) ToUsecase(c fiber.Ctx) policyUC.AddPermissionRequest {
	return policyUC.AddPermissionRequest{
		Name:              e.Name,
		Description:       e.Description,
		AssignToResources: e.AssignToResources,
	}
}

type AddRoleRequest struct {
	Name                string  `json:"name" validate:"required"`
	Description         string  `json:"description" validate:"required"`
	AssignToPermissions []int64 `json:"assign_to_permissions"`
}

func (e *AddRoleRequest) ToUsecase(c fiber.Ctx) policyUC.AddRoleRequest {
	return policyUC.AddRoleRequest{
		Name:                e.Name,
		Description:         e.Description,
		AssignToPermissions: e.AssignToPermissions,
	}
}

type UpdateResourceRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Type   string `json:"type" validate:"oneof=http grpc"`
	Method string `json:"method"`
	Path   string `json:"path" validate:"required"`
}

func (e *UpdateResourceRequest) ToUsecase(c fiber.Ctx) policyUC.UpdateResourceRequest {
	return policyUC.UpdateResourceRequest{
		ID:     e.ID,
		Type:   policyRepo.ServiceType(e.Type),
		Method: e.Method,
		Path:   e.Path,
	}
}

type UpdatePermissionRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (e *UpdatePermissionRequest) ToUsecase(c fiber.Ctx) policyUC.UpdatePermissionRequest {
	return policyUC.UpdatePermissionRequest{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}

type UpdateRoleRequest struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (e *UpdateRoleRequest) ToUsecase(c fiber.Ctx) policyUC.UpdateRoleRequest {
	return policyUC.UpdateRoleRequest{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}

type GetUserPermissionsResponse []string

func (e *GetUserPermissionsResponse) FromUsecase(resp policyUC.GetUserPermissionsResponse) {
	if e == nil {
		*e = make([]string, len(resp.Permissions))
	}
	for _, v := range resp.Permissions {
		*e = append(*e, v.Name)
	}
}
