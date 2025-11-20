package http

import (
	"context"
	"fmt"

	"github.com/aidapedia/gdk/http/server"
	gmask "github.com/aidapedia/gdk/mask"
	"github.com/gofiber/fiber/v3"

	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/internal/interface/http/middleware"
	"github.com/aidapedia/jabberwock/internal/interface/http/route"
	"github.com/aidapedia/jabberwock/pkg/config"
)

// HTTPServiceInterface is an interface to handle http service
type HTTPServiceInterface interface {
	ListenAndServe() error
}

// HTTPService is a struct to handle http service
type HTTPService struct {
	svr *server.Server
}

// NewHTTPService is a function to create a new http service
func NewHTTPService(handler *handler.Handler, middleware *middleware.Middleware) HTTPServiceInterface {
	cfg := config.GetConfig(context.Background())
	svr, _ := server.NewWithDefaultConfig(cfg.App.Name, gmask.NewDefault())
	route.Register(svr.App, handler, middleware)

	return &HTTPService{
		svr: svr,
	}
}

// ListenAndServe is a function to start http service
func (h *HTTPService) ListenAndServe() error {
	cfg := config.GetConfig(context.Background())
	fmt.Printf("Starting HTTP server %s:%d\n", cfg.App.HTTPServer.Address, cfg.App.HTTPServer.Port)
	return h.svr.Listen(fmt.Sprintf("%s:%d", cfg.App.HTTPServer.Address, cfg.App.HTTPServer.Port), fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}
