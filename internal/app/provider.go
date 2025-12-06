package app

import (
	"github.com/aidapedia/jabberwock/internal/app/service"

	httpInterface "github.com/aidapedia/jabberwock/internal/interface/http"
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/internal/interface/http/middleware"

	"github.com/aidapedia/jabberwock/internal/driver/casbin"
	"github.com/aidapedia/jabberwock/internal/driver/database"
	"github.com/aidapedia/jabberwock/internal/driver/redis"

	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"

	authenticatedUsecase "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	policyUsecase "github.com/aidapedia/jabberwock/internal/usecase/policy"
	userdatacenterUsecase "github.com/aidapedia/jabberwock/internal/usecase/userdatacenter"

	"github.com/google/wire"
	_ "github.com/lib/pq"
)

var (
	driverSet = wire.NewSet(
		redis.NewRedis,
		casbin.NewCasbin,
		database.NewDatabase,
	)

	repositorySet = wire.NewSet(
		sessionRepo.New,
		userRepo.New,
		policyRepo.New,
	)

	usecaseSet = wire.NewSet(
		authenticatedUsecase.New,
		policyUsecase.New,
		userdatacenterUsecase.New,
	)

	httpSet = wire.NewSet(
		driverSet,
		repositorySet,
		usecaseSet,

		middleware.NewMiddleware,
		handler.NewHandler,
		httpInterface.NewHTTPService,

		service.NewServiceHTTP,
	)
)
