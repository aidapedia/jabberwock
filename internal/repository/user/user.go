package user

import (
	"context"
	"time"

	"github.com/aidapedia/gdk/telemetry/tracer"
)

func (r *Repository) FindByID(ctx context.Context, id int64) (resp User, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserRepository/FindByID")
	defer span.Finish(err)

	query := queryFindByID
	args := []interface{}{id}
	err = r.database.QueryRowContext(ctx, query, args...).
		Scan(&resp.ID, &resp.Name, &resp.Phone, &resp.Email, &resp.Password,
			&resp.Type, &resp.Status, &resp.IsVerified, &resp.AvatarURL, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repository) FindByPhone(ctx context.Context, phone string) (resp User, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserRepository/FindByPhone")
	defer span.Finish(err)

	query := queryFindByPhone
	args := []interface{}{phone}
	err = r.database.QueryRowContext(ctx, query, args...).
		Scan(&resp.ID, &resp.Name, &resp.Phone, &resp.Email, &resp.Password,
			&resp.Type, &resp.Status, &resp.IsVerified, &resp.AvatarURL, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (resp User, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserRepository/FindByEmail")
	defer span.Finish(err)

	query := queryFindByEmail
	args := []interface{}{email}
	err = r.database.QueryRowContext(ctx, query, args...).
		Scan(&resp.ID, &resp.Name, &resp.Phone, &resp.Email, &resp.Password,
			&resp.Type, &resp.Status, &resp.IsVerified, &resp.AvatarURL, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int64, status Status) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserRepository/UpdateStatus")
	defer span.Finish(err)

	query := queryUpdateStatus
	args := []interface{}{status, time.Now(), id}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user *User) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserRepository/CreateUser")
	defer span.Finish(err)

	query := queryCreateUser
	args := []interface{}{user.Name, user.Phone, user.Email, user.Password, user.Type,
		user.Status, user.IsVerified, user.AvatarURL, time.Now(), time.Now()}
	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
