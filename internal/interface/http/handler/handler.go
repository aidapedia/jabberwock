package handler

import (
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/aidapedia/jabberwock/internal/usecase/userdatacenter"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	authUsecase authenticated.Interface
	userUsecase userdatacenter.Interface
}

func NewHandler(authUsecase authenticated.Interface, userUsecase userdatacenter.Interface) *Handler {
	return &Handler{authUsecase: authUsecase, userUsecase: userUsecase}
}

func (h *Handler) Ping(c fiber.Ctx) error {
	ghttp.JSONResponse(c, map[string]string{"message": "pong"}, nil)
	return nil
}
