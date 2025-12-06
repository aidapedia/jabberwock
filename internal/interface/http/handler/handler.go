package handler

import (
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/aidapedia/jabberwock/internal/usecase/policy"
	"github.com/aidapedia/jabberwock/internal/usecase/userdatacenter"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	authUsecase   authenticated.Interface
	policyUsecase policy.Interface
	userUsecase   userdatacenter.Interface
}

func NewHandler(authUsecase authenticated.Interface, policyUsecase policy.Interface, userUsecase userdatacenter.Interface) *Handler {
	return &Handler{
		authUsecase:   authUsecase,
		policyUsecase: policyUsecase,
		userUsecase:   userUsecase,
	}
}

func (h *Handler) Ping(c fiber.Ctx) error {
	ghttp.JSONResponse(c, &ghttp.SuccessResponse{
		Data: map[string]string{"message": "pong"},
	}, nil)
	return nil
}
