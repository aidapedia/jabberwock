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

func (r *Repository) SetActiveSession(ctx context.Context, session Session) error {
	return nil
}

func (r *Repository) FindActiveSessionByTokenID(ctx context.Context, tokenID string) (Session, error) {
	return Session{}, nil
}

func (r *Repository) DeleteActiveSession(ctx context.Context, tokenID string) error {
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
