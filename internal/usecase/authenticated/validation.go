package authenticated

import (
	"context"
	"errors"
	"net/http"

	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
	pkgLog "github.com/aidapedia/jabberwock/pkg/log"

	gcryptography "github.com/aidapedia/gdk/cryptography"
	gers "github.com/aidapedia/gdk/error"
)

// validationPassword is a function to validate password
func (a *Usecase) validationPassword(ctx context.Context, user userRepo.User, password string) error {
	// Password check attempt
	// will check if user has attempt login before
	lastAttempt, err := a.checkAttemptFailed(ctx, user)
	if err != nil {
		return err
	}
	// And then if no blocked check password is valid or not.
	if !gcryptography.CheckHash(password, user.Password) {
		// Not valid will update attempt login and return error
		err = a.updateAttemptFailed(ctx, user, lastAttempt)
		if err != nil {
			return err
		}
		return gers.NewWithMetadata(errors.New("password is incorrect"), pkgLog.Metadata(http.StatusBadRequest, "Password is incorrect"))
	}
	// Success login will reset attempt login
	err = a.resetAttemptFailed(ctx, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// validateUser is a function to validate user
func (a *Usecase) validateUser(user userRepo.User) error {
	if user.Status == userRepo.StatusBlocked {
		return gers.NewWithMetadata(errors.New("account is blocked"),
			pkgLog.Metadata(http.StatusBadRequest, "Your account is blocked"))
	}
	if user.IsVerified == userRepo.VerifiedNone {
		return gers.NewWithMetadata(errors.New("account is not verified"),
			pkgLog.Metadata(http.StatusBadRequest, "Please verify your identity first"))
	}
	return nil
}
