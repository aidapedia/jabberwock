package userdatacenter

import (
	"context"
	"net/http"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"

	cerror "github.com/aidapedia/jabberwock/internal/common/error"
)

func (uc *Usecase) GetUserByID(ctx context.Context, id int64) (resp User, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserDataCenterUsecase/GetUserByID")
	defer span.Finish(err)

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if err == cerror.ErrorNotFound {
			return resp, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusBadRequest, "User not found"))
		}
		return resp, gers.NewWithMetadata(err, ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	resp.Transform(user)
	return resp, nil
}
