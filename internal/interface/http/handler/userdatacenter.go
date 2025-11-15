package handler

import (
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/gdk/util"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) GetUserByID(c fiber.Ctx) error {
	span, ctx := tracer.StartSpanFromContext(c.Context(), "UserDataCenterHandler/GetUserByID")
	defer span.Finish(nil)

	resp, err := h.userUsecase.GetUserByID(ctx, util.ToInt64(c.Params("id")))
	if err != nil {
		ghttp.JSONResponse(c, nil, err)
		return err
	}

	ghttp.JSONResponse(c, &ghttp.SuccessResponse{
		Data: resp,
	}, nil)
	return nil
}
