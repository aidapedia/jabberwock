package middleware

import (
	"context"

	"github.com/aidapedia/gdk/log"
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	authUcesace "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"go.uber.org/zap"
)

type Middleware struct {
	authenticatedUC authUcesace.Interface
}

func NewMiddleware(authenticatedUC authUcesace.Interface) *Middleware {
	ctx := context.Background()
	err := authenticatedUC.LoadPolicy(ctx, policyRepo.HTTPServiceType)
	if err != nil {
		log.FatalCtx(ctx, "Error when setting policy", zap.Error(err))
	}
	return &Middleware{authenticatedUC: authenticatedUC}
}
