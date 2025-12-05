package route

import (
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/aidapedia/jabberwock/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, handler *handler.Handler, middleware *middleware.Middleware) {
	// Public Routes
	app.Get("/ping", handler.Ping)
	app.Group("/auth").
		Post("/login", handler.Login).
		Post("/register", handler.Register)

	// Protected Routes
	protected := app.Group("", middleware.CheckAccess())
	protected.Group("/auth").
		Post("/logout", handler.Logout)
	protected.Group("/policy").
		Post("/resource", handler.AddResource).
		Post("/permission", handler.AddPermission).
		Post("/role", handler.AddRole).
		Put("/resource", handler.UpdateResource).
		Put("/permission", handler.UpdatePermission).
		Put("/role", handler.UpdateRole).
		Delete("/resource/:id", handler.DeleteResource).
		Delete("/permission/:id", handler.DeletePermission).
		Delete("/role/:id", handler.DeleteRole)
	protected.Get("/user/:id", handler.GetUserByID)
}
