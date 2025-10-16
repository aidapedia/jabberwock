package session

import "context"

func (r *Repository) SetActiveSession(ctx context.Context, session Session) error {
	return nil
}

func (r *Repository) FindActiveSessionByTokenID(ctx context.Context, tokenID string) (Session, error) {
	return Session{}, nil
}

func (r *Repository) DeleteActiveSession(ctx context.Context, tokenID string) error {
	return nil
}

func (r *Repository) GetLoginAttempt(ctx context.Context, userID int64) (LoginAttempt, error) {
	return LoginAttempt{}, nil
}

func (r *Repository) SetLoginAttempt(ctx context.Context, userID int64, loginAttempt LoginAttempt) error {
	return nil
}
