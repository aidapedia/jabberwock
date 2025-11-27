package middleware

import (
	ghttp "github.com/aidapedia/gdk/http"
	gmiddleware "github.com/aidapedia/gdk/http/server/middleware"

	authUsecase "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
	"github.com/gofiber/fiber/v3"
)

// CheckAccess checks if the request has a valid access token
func (e *Middleware) CheckAccess() gmiddleware.Middleware {
	return func(c fiber.Ctx) error {
		err := e.authenticatedUC.CheckAccessToken(c.Context(), authUsecase.CheckAccessTokenPayload{
			Token: string(c.Request().Header.Peek("Authorization")),
			// ElementID is the path of the request
			// e.g. GET|/v1/users
			ElementID: authUsecase.GenerateElementID(c.Method(), c.Path()),
		})
		if err != nil {
			ghttp.JSONResponse(c, nil, err)
			return err
		}

		// Since the user id session will get once each request, we can add the user data to the context
		// c.Locals(commonCtx.ContextKeyUserID, userID)
		return c.Next()
	}
}
