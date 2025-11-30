package policy

import (
	"context"
	"database/sql"
)

type Interface interface {
	// CreateRole(ctx context.Context, role Role) (err error)
	// CreatePermission(ctx context.Context, permission Permission) (err error)
	// CreateResource(ctx context.Context, resource Resource) (err error)
	// CreateRolePermission(ctx context.Context, roleID, permissionID int64) (err error)
	// CreateResourcePermission(ctx context.Context, resourceID, permissionID int64) (err error)

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
