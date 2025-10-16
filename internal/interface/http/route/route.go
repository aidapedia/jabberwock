package route

import (
	"github.com/aidapedia/jabberwock/internal/interface/http/handler"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, handler *handler.Handler) {
	app.Post("/login", handler.Login)
	app.Get("/ping", handler.Ping)
}
