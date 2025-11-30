package middleware

import (
	authUcesace "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
)

type Middleware struct {
	authenticatedUC authUcesace.Interface
}

func NewMiddleware(authenticatedUC authUcesace.Interface) *Middleware {
	return &Middleware{authenticatedUC: authenticatedUC}
}
