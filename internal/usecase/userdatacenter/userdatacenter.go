package userdatacenter

import (
	"context"

	"github.com/aidapedia/gdk/telemetry/tracer"
)

func (uc *Usecase) GetUserByID(ctx context.Context, id int64) (resp User, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UserDataCenterUsecase/GetUserByID")
	defer span.Finish(err)

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return resp, err
	}

	resp.Transform(user)
	return resp, nil
}
