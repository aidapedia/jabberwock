package route

import (
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, handler *handler.Handler, middleware *middleware.Middleware) {
	// Public Routes
	app.Get("/ping", handler.Ping)
	app.Post("/login", handler.Login)
	app.Post("/register", handler.Register)

	// Protected Routes
	auth := app.Group("", middleware.CheckAccess())
	auth.Post("/logout", handler.Logout)
}
