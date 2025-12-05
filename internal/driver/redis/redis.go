package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/aidapedia/gdk/cache/engine"
	gredisengine "github.com/aidapedia/gdk/cache/engine"
	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/jabberwock/pkg/config"
	goredis "github.com/redis/go-redis/v9"
	goredismaint "github.com/redis/go-redis/v9/maintnotifications"
	"go.uber.org/zap"
)

func NewRedis(ctx context.Context) engine.Interface {
	cfg := config.GetConfig(ctx)
	if cfg == nil {
		log.FatalCtx(ctx, "failed to connect redis: %v", zap.Error(errors.New("config is nil")))
	}

	redis, err := gredisengine.NewGoRedisClient(gredisengine.GoRedisClientOpt{
		Opt: &goredis.Options{
			Addr: fmt.Sprintf("%s:%d", cfg.Secret.Redis.Address, cfg.Secret.Redis.Port),
			MaintNotificationsConfig: &goredismaint.Config{
				Mode: goredismaint.ModeDisabled,
			},
		},
	})
	if err != nil {
		log.FatalCtx(ctx, "failed to connect redis: %v", zap.Error(err))
	}

	return redis
}
