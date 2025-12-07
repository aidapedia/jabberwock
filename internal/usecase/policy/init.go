package policy

import (
	"context"

	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	"github.com/casbin/casbin/v2"
)

type Interface interface {
	// Policy Management
	LoadPolicy(ctx context.Context, serviceType policyRepo.ServiceType) (err error)
	AddResource(ctx context.Context, req AddResourceRequest) (err error)
	AddPermission(ctx context.Context, req AddPermissionRequest) (err error)
	AddRole(ctx context.Context, req AddRoleRequest) (err error)
	UpdateResource(ctx context.Context, req UpdateResourceRequest) (err error)
	UpdatePermission(ctx context.Context, req UpdatePermissionRequest) (err error)
	UpdateRole(ctx context.Context, req UpdateRoleRequest) (err error)
	DeleteResource(ctx context.Context, req DeleteResourceRequest) (err error)
	DeletePermission(ctx context.Context, req DeletePermissionRequest) (err error)
	DeleteRole(ctx context.Context, req DeleteRoleRequest) (err error)
	GetUserPermissions(ctx context.Context, userID int64) (resp GetUserPermissionsResponse, err error)
}

type Usecase struct {
	enforcer   *casbin.Enforcer
	policyRepo policyRepo.Interface
}

func New(policyRepo policyRepo.Interface, enforcer *casbin.Enforcer) Interface {
	return &Usecase{
		enforcer:   enforcer,
		policyRepo: policyRepo,
	}
}
