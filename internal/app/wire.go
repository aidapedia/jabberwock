//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package app

import (
	"github.com/aidapedia/jabberwock/internal/interface/http"
	"github.com/google/wire"
)

func InitHTTPServer() http.HTTPServiceInterface {
	wire.Build(httpSet)
	return &http.HTTPService{}
}
