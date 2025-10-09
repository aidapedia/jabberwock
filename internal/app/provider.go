package app

import (
	"errors"
	"fmt"
	"log"

	httpInterface "github.com/aidapedia/jabberwock/internal/interface/http"
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/pkg/config"

	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"

	authenticatedUsecase "github.com/aidapedia/jabberwock/internal/usecase/authenticated"

	gredisengine "github.com/aidapedia/gdk/cache/engine"
	casbin "github.com/casbin/casbin/v2"
	casbinUtil "github.com/casbin/casbin/v2/util"
	goredis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func redisProvider() gredisengine.Interface {
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("failed to connect redis: %v", errors.New("config is nil"))
	}

	redis, err := gredisengine.NewGoRedisClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Storage.Redis.Address, cfg.Storage.Redis.Port),
	})
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	return redis
}

// ProviderCasbin is a function to create a new casbin enforcer
func casbinProvider() *casbin.Enforcer {
	cfg := config.GetConfig()
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
		handler.NewHandler,
		httpInterface.NewHTTPService,
	)
)
