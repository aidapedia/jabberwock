package policy

import (
	"context"
	"database/sql"
)

type Interface interface {
	CreateResource(ctx context.Context, resource Resource) (err error)
	BulkAssignResources(ctx context.Context, tx *sql.Tx, permissionID int64, resourceIDs []int64) (err error)
	UpdateResource(ctx context.Context, resource Resource) (err error)
	DeleteResource(ctx context.Context, resourceID int64) (err error)
	BulkDeleteAssignResource(ctx context.Context, tx *sql.Tx, permissionID int64, resourceIDs []int64) (err error)

	CreatePermission(ctx context.Context, permission Permission) (err error)
	BulkAssignPermissions(ctx context.Context, tx *sql.Tx, roleID int64, permissionIDs []int64) (err error)
	UpdatePermission(ctx context.Context, permission Permission) (err error)
	DeletePermission(ctx context.Context, permissionID int64) (err error)
	BulkDeleteAssignPermission(ctx context.Context, tx *sql.Tx, roleID int64, permissionIDs []int64) (err error)

	CreateRole(ctx context.Context, role Role) (err error)
	AssignRole(ctx context.Context, userID int64, roleID int64) (err error)
	GetRoleByUserID(ctx context.Context, userID int64) (resp GetRoleByUserIDResp, err error)
	UpdateRole(ctx context.Context, role Role) (err error)
	DeleteRole(ctx context.Context, roleID int64) (err error)

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
