package handler

import (
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	usecase authenticated.Interface
}

func NewHandler(usecase authenticated.Interface) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Ping(c fiber.Ctx) error {
	ghttp.JSONResponse(c, map[string]string{"message": "pong"}, nil)
	return nil
}
