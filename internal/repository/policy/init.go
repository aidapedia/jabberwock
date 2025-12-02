package policy

import (
	"context"
	"database/sql"
)

type Interface interface {
	CreateResource(ctx context.Context, resource Resource) (err error)
	BulkAssignResources(ctx context.Context, permissionID int64, resourceID []int64) (err error)

	CreatePermission(ctx context.Context, permission Permission) (err error)
	BulkAssignPermissions(ctx context.Context, roleID int64, permissionID []int64) (err error)

	CreateRole(ctx context.Context, role Role) (err error)
	AssignRole(ctx context.Context, userID int64, roleID int64) (err error)
	GetRoleByUserID(ctx context.Context, userID int64) (resp GetRoleByUserIDResp, err error)

	LoadPolicy(ctx context.Context, serviceType ServiceType) (resp LoadPolicyResponse, err error)
}

type Repository struct {
	database *sql.DB
}

func New(database *sql.DB) Interface {
	return &Repository{
		database: database,
	}
}
