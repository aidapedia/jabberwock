package userdatacenter

import "context"

func (uc *Usecase) GetUserByID(ctx context.Context, id int64) (User, error) {
	return User{}, nil
}
