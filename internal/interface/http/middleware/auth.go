package middleware

import (
	ghttp "github.com/aidapedia/gdk/http"
	gmiddleware "github.com/aidapedia/gdk/http/server/middleware"

	"github.com/aidapedia/jabberwock/internal/common/constant"
	authUsecase "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/gofiber/fiber/v3"
)

// CheckAccess checks if the request has a valid access token
func (e *Middleware) CheckAccess() gmiddleware.Middleware {
	return func(c fiber.Ctx) error {
		resp, err := e.authenticatedUC.CheckAccessToken(c.Context(), authUsecase.CheckAccessTokenPayload{
			Token:     string(c.Request().Header.Peek("Authorization")),
			ElementID: authUsecase.GenerateElementID(c.Method(), c.Path()),
		})
		if err != nil {
			ghttp.JSONResponse(c, nil, err)
			return err
		}

		c.Locals(constant.ContextKeyUserID, resp.UserID)
		return c.Next()
	}
}
