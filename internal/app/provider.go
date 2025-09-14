package app

import (
	"github.com/kurniajigunawan/homestay/internal/app/service"

	httpInterface "github.com/kurniajigunawan/homestay/internal/interface/http"
	"github.com/kurniajigunawan/homestay/internal/interface/http/handler"
	"github.com/kurniajigunawan/homestay/internal/interface/http/middleware"

	"github.com/kurniajigunawan/homestay/internal/driver/casbin"
	"github.com/kurniajigunawan/homestay/internal/driver/database"
	"github.com/kurniajigunawan/homestay/internal/driver/redis"

	policyRepo "github.com/kurniajigunawan/homestay/internal/repository/policy"
	sessionRepo "github.com/kurniajigunawan/homestay/internal/repository/session"
	userRepo "github.com/kurniajigunawan/homestay/internal/repository/user"

	authenticatedUsecase "github.com/kurniajigunawan/homestay/internal/usecase/authenticated"
	policyUsecase "github.com/kurniajigunawan/homestay/internal/usecase/policy"
	userdatacenterUsecase "github.com/kurniajigunawan/homestay/internal/usecase/userdatacenter"

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
