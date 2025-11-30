package policy

import (
	"context"
	"database/sql"

	"github.com/aidapedia/gdk/telemetry/tracer"
)

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
