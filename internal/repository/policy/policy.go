package policy

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/aidapedia/gdk/telemetry/tracer"
)

func (r *Repository) CreateRole(ctx context.Context, role Role) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/CreateRole")
	defer span.Finish(err)

	query := queryCreateRole
	args := []interface{}{role.Name, role.Description}
	err = r.database.QueryRowContext(ctx, query, args...).Scan(&role.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AssignRole(ctx context.Context, userID int64, roleID int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/AssignRole")
	defer span.Finish(err)

	query := queryAssignRole
	args := []interface{}{userID, roleID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetRoleByUserID(ctx context.Context, userID int64) (resp GetRoleByUserIDResp, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/GetRoleByUserID")
	defer span.Finish(err)

	query := queryGetRoleByUserID
	args := []interface{}{userID}
	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var role Role
		err = rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return []Role{}, err
		}
		resp = append(resp, role)
	}

	return resp, nil
}

func (r *Repository) LoadPolicy(ctx context.Context, serviceType ServiceType) (resp LoadPolicyResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/LoadPolicy")
	defer span.Finish(err)

	query := queryLoadPolicy
	args := []interface{}{serviceType}
	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var policy Policy
		err = rows.Scan(&policy.Role, &policy.Type, &policy.Path, &policy.Method)
		if err != nil {
			return []Policy{}, err
		}
		resp = append(resp, policy)
	}

	return resp, nil
}

func (r *Repository) CreateResource(ctx context.Context, resource Resource) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/CreateResource")
	defer span.Finish(err)

	query := queryCreateResource
	args := []interface{}{resource.Type, resource.Method, resource.Path}
	err = r.database.QueryRowContext(ctx, query, args...).Scan(&resource.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) BulkAssignResources(ctx context.Context, tx *sql.Tx, permissionID int64, resourceIDs []int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/BulkAssignResources")
	defer span.Finish(err)

	query := queryBulkAssignResource
	args := []interface{}{}

	var valCount int
	var values []string
	for _, resourceID := range resourceIDs {
		values = append(values, fmt.Sprintf("($%d, $%d)", valCount+1, valCount+2))
		args = append(args, permissionID, resourceID)
		valCount += 2
	}
	query += fmt.Sprintf("VALUES %s;", strings.Join(values, ","))

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.database.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreatePermission(ctx context.Context, permission Permission) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/CreatePermission")
	defer span.Finish(err)

	query := queryCreatePermission
	args := []interface{}{permission.Name, permission.Description}
	err = r.database.QueryRowContext(ctx, query, args...).Scan(&permission.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) BulkAssignPermissions(ctx context.Context, tx *sql.Tx, roleID int64, permissionIDs []int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/BulkAssignPermissions")
	defer span.Finish(err)

	query := queryBulkAssignPermission
	args := []interface{}{}

	var (
		valCount int
		values   []string
	)
	for _, permissionID := range permissionIDs {
		values = append(values, fmt.Sprintf("($%d, $%d)", valCount+1, valCount+2))
		args = append(args, roleID, permissionID)
		valCount += 2
	}
	query += fmt.Sprintf("VALUES %s;", strings.Join(values, ","))

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.database.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) BulkDeleteAssignPermission(ctx context.Context, tx *sql.Tx, roleID int64, permissionIDs []int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/BulkDeleteAssignPermission")
	defer span.Finish(err)

	query := queryBulkDeleteAssignPermission
	args := []interface{}{}

	var (
		valCount int
		values   []string
	)
	for _, permissionID := range permissionIDs {
		values = append(values, fmt.Sprintf("($%d, $%d)", valCount+1, valCount+2))
		args = append(args, roleID, permissionID)
		valCount += 2
	}
	query += fmt.Sprintf("VALUES %s;", strings.Join(values, ","))

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.database.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdatePermission(ctx context.Context, permission Permission) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/UpdatePermission")
	defer span.Finish(err)

	query := queryUpdatePermission
	args := []interface{}{permission.Name, permission.Description, permission.ID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeletePermission(ctx context.Context, permissionID int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/DeletePermission")
	defer span.Finish(err)

	query := queryDeletePermission
	args := []interface{}{permissionID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) BulkDeleteAssignResource(ctx context.Context, tx *sql.Tx, permissionID int64, resourceIDs []int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/BulkDeleteAssignResource")
	defer span.Finish(err)

	query := queryBulkDeleteAssignResource
	args := []interface{}{}

	var (
		valCount int
		values   []string
	)
	for _, resourceID := range resourceIDs {
		values = append(values, fmt.Sprintf("($%d, $%d)", valCount+1, valCount+2))
		args = append(args, permissionID, resourceID)
		valCount += 2
	}
	query += fmt.Sprintf("VALUES %s;", strings.Join(values, ","))

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.database.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateResource(ctx context.Context, resource Resource) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/UpdateResource")
	defer span.Finish(err)

	query := queryUpdateResource
	args := []interface{}{resource.Type, resource.Method, resource.Path, resource.ID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteResource(ctx context.Context, resourceID int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/DeleteResource")
	defer span.Finish(err)

	query := queryDeleteResource
	args := []interface{}{resourceID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateRole(ctx context.Context, role Role) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/UpdateRole")
	defer span.Finish(err)

	query := queryUpdateRole
	args := []interface{}{role.Name, role.Description, role.ID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRole(ctx context.Context, roleID int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/DeleteRole")
	defer span.Finish(err)

	query := queryDeleteRole
	args := []interface{}{roleID}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserPermissions(ctx context.Context, userID int64) (resp []Permission, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/GetUserPermissions")
	defer span.Finish(err)

	query := queryGetUserPermissions
	args := []interface{}{userID}
	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission Permission
		err = rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		resp = append(resp, permission)
	}
	return resp, nil
}

func (r *Repository) GetAllPermissions(ctx context.Context) (resp []Permission, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/GetAllPermissions")
	defer span.Finish(err)

	query := queryGetAllPermissions
	rows, err := r.database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission Permission
		err = rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		resp = append(resp, permission)
	}
	return resp, nil
}
