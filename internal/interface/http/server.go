package http

import (
	"fmt"

	"github.com/aidapedia/gdk/http/server"
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/internal/interface/http/route"
	"github.com/aidapedia/jabberwock/pkg/config"
	"github.com/gofiber/fiber/v3"
)

// HTTPServiceInterface is an interface to handle http service
type HTTPServiceInterface interface {
	ListenAndServe() error
	// Stop()
}

// HTTPService is a struct to handle http service
type HTTPService struct {
	svr *server.Server
}

// NewHTTPService is a function to create a new http service
func NewHTTPService(handler *handler.Handler) HTTPServiceInterface {
	svr, _ := server.NewWithDefaultConfig(server.WithMiddlewares(func(c fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	}))
	route.Register(svr.GetFiberApp(), handler)

	return &HTTPService{
		svr: svr,
	}
}

// ListenAndServe is a function to start http service
func (h *HTTPService) ListenAndServe() error {
	cfg := config.GetConfig()
	return h.svr.ListenGracefully(fmt.Sprintf("%s:%d", cfg.App.HTTPServer.Address, cfg.App.HTTPServer.Port))
}
