package user

import "context"

func (r *Repository) FindByID(ctx context.Context, id int64) (User, error) {
	return User{}, nil
}

func (r *Repository) FindByPhone(ctx context.Context, phone string) (User, error) {
	return User{}, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (User, error) {
	return User{}, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int64, status Status) error {
	return nil
}
