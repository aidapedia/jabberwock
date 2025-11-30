//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package app

import (
	"context"

	"github.com/aidapedia/jabberwock/internal/app/service"
	"github.com/google/wire"
)

func InitHTTPServer(ctx context.Context) *service.ServiceHTTP {
	wire.Build(httpSet)
	return &service.ServiceHTTP{}
}
