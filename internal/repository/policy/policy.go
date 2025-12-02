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
		err = rows.Scan(&role.Name, &role.Description)
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

func (r *Repository) BulkAssignResources(ctx context.Context, permissionID int64, resourceID []int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/BulkAssignResources")
	defer span.Finish(err)

	query := queryBulkAssignResource
	args := []interface{}{}

	var valCount int
	var values []string
	for _, resourceID := range resourceID {
		values = append(values, fmt.Sprintf("($%d, $%d)", valCount+1, valCount+2))
		args = append(args, permissionID, resourceID)
		valCount += 2
	}
	query += fmt.Sprintf("VALUES %s;", strings.Join(values, ","))

	_, err = r.database.ExecContext(ctx, query, args...)
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

func (r *Repository) BulkAssignPermissions(ctx context.Context, roleID int64, permissionID []int64) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "PolicyRepository/BulkAssignPermissions")
	defer span.Finish(err)

	query := queryBulkAssignPermission
	args := []interface{}{}

	var (
		valCount int
		values   []string
	)
	for _, permissionID := range permissionID {
		values = append(values, fmt.Sprintf("($%d, $%d)", valCount+1, valCount+2))
		args = append(args, roleID, permissionID)
		valCount += 2
	}
	query += fmt.Sprintf("VALUES %s;", strings.Join(values, ","))

	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
