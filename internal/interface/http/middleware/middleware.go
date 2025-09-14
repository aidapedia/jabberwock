package middleware

import authUcesace "github.com/aidapedia/jabberwock/internal/usecase/authenticated"

type Middleware struct {
	authenticatedUC authUcesace.Interface
}
