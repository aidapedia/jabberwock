package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/jabberwock/internal/common/cache"
	"github.com/redis/go-redis/v9"
)

func (r *Repository) CreateActiveSession(ctx context.Context, session *Session) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "SessionRepository/CreateActiveSession")
	defer span.Finish(err)

	if session == nil {
		err = fmt.Errorf("insert data request is nil")
		return err
	}

	query := queryCreateActiveSession
	args := []interface{}{session.Token, session.UserID, session.UserAgent, session.IP, time.Now(), time.Now()}

	err = r.database.QueryRowContext(ctx, query, args...).Scan(&session.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindActiveSessionByTokenID(ctx context.Context, tokenID string) (resp Session, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "SessionRepository/FindActiveSessionByTokenID")
	defer span.Finish(err)

	query := queryFindActiveSessionByTokenID
	args := []interface{}{tokenID}
	err = r.database.QueryRowContext(ctx, query, args...).
		Scan(&resp.ID, &resp.UserID, &resp.UserAgent, &resp.IP, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *Repository) DeleteActiveSession(ctx context.Context, tokenID string) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "SessionRepository/DeleteActiveSession")
	defer span.Finish(err)

	query := queryDeleteActiveSessionByTokenID
	args := []interface{}{tokenID}

	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLoginAttempt(ctx context.Context, userID int64) (result LoginAttempt, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "SessionRepository/GetLoginAttempt")
	defer span.Finish(err)

	key := fmt.Sprintf(cache.RedisKeyLoginAttempt, userID)
	val, err := r.redis.GET(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return result, err
		}
		return LoginAttempt{}, err
	}

	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *Repository) SetLoginAttempt(ctx context.Context, userID int64, loginAttempt LoginAttempt) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "SessionRepository/SetLoginAttempt")
	defer span.Finish(err)

	key := fmt.Sprintf(cache.RedisKeyLoginAttempt, userID)
	return r.redis.SET(ctx, key, loginAttempt, time.Until(loginAttempt.RefreshTime))
}
