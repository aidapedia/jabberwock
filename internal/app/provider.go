package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	httpInterface "github.com/aidapedia/jabberwock/internal/interface/http"
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/internal/interface/http/middleware"
	"github.com/aidapedia/jabberwock/pkg/config"

	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"

	authenticatedUsecase "github.com/aidapedia/jabberwock/internal/usecase/authenticated"

	gredisengine "github.com/aidapedia/gdk/cache/engine"
	casbin "github.com/casbin/casbin/v2"
	casbinUtil "github.com/casbin/casbin/v2/util"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	goredis "github.com/redis/go-redis/v9"
	goredismaint "github.com/redis/go-redis/v9/maintnotifications"
)

var DatabaseDriver *sql.DB

func databaseProvider(ctx context.Context) *sql.DB {
	cfg := config.GetConfig(ctx)
	if cfg == nil {
		log.Fatalf("failed to connect database: %v", errors.New("config is nil"))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Secret.Database.Host, cfg.Secret.Database.Port, cfg.Secret.Database.Username, cfg.Secret.Database.Password, cfg.Secret.Database.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	DatabaseDriver = db
	return db
}

func redisProvider(ctx context.Context) gredisengine.Interface {
	cfg := config.GetConfig(ctx)
	if cfg == nil {
		log.Fatalf("failed to connect redis: %v", errors.New("config is nil"))
	}

	redis, err := gredisengine.NewGoRedisClient(gredisengine.GoRedisClientOpt{
		Opt: &goredis.Options{
			Addr: fmt.Sprintf("%s:%d", cfg.Storage.Redis.Address, cfg.Storage.Redis.Port),
			MaintNotificationsConfig: &goredismaint.Config{
				Mode: goredismaint.ModeDisabled,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	return redis
}

// ProviderCasbin is a function to create a new casbin enforcer
func casbinProvider(ctx context.Context) *casbin.Enforcer {
	cfg := config.GetConfig(ctx)
	authEnforcer, err := casbin.NewEnforcer(cfg.App.Auth.ModelPath, cfg.App.Auth.PolicyPath)
	if err != nil {
		log.Fatal(err)
	}
	authEnforcer.AddNamedMatchingFunc("g", "KeyMatch2", casbinUtil.KeyMatch2)
	return authEnforcer
}

var (
	driverSet = wire.NewSet(
		redisProvider,
		casbinProvider,
		databaseProvider,
	)

	repositorySet = wire.NewSet(
		sessionRepo.New,
		userRepo.New,
	)

	usecaseSet = wire.NewSet(
		authenticatedUsecase.New,
	)

	httpSet = wire.NewSet(
		driverSet,
		repositorySet,
		usecaseSet,
		middleware.NewMiddleware,
		handler.NewHandler,
		httpInterface.NewHTTPService,
	)
)
