package policy

import (
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
)

type AddResourceRequest struct {
	Type   policyRepo.ServiceType
	Method string
	Path   string
}

type AddPermissionRequest struct {
	Name              string
	Description       string
	AssignToResources []int64
}

type AddRoleRequest struct {
	Name                string
	Description         string
	AssignToPermissions []int64
}

type UpdateResourceRequest struct {
	ID     int64
	Type   policyRepo.ServiceType
	Method string
	Path   string
}

type UpdatePermissionRequest struct {
	ID          int64
	Name        string
	Description string
}

type UpdateRoleRequest struct {
	ID          int64
	Name        string
	Description string
}

type DeleteResourceRequest struct {
	ID int64
}

type DeletePermissionRequest struct {
	ID int64
}

type DeleteRoleRequest struct {
	ID int64
}
