package authenticated

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	gers "github.com/aidapedia/gdk/error"
	constant "github.com/aidapedia/jabberwock/internal/common/constant"
	sessionRepo "github.com/aidapedia/jabberwock/internal/repository/session"
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
	pkgLog "github.com/aidapedia/jabberwock/pkg/log"
)

var attemptConfig = map[int]time.Duration{
	0: 0,
	1: 0,
	2: 30 * time.Second,
	3: 1 * time.Minute,
	4: 15 * time.Minute,
	5: 1 * time.Hour,
	6: 6 * time.Hour,
	7: 24 * time.Hour,
}

func (a *Usecase) checkAttemptFailed(ctx context.Context, user userRepo.User) (int, error) {
	loginResp, err := a.sessionRepo.GetLoginAttempt(ctx, user.ID)
	if err != nil {
		return -1, gers.NewWithMetadata(err,
			pkgLog.Metadata(http.StatusInternalServerError, constant.ErrorMessageTryAgain))
	}
	if loginResp.Attempt == 0 {
		return loginResp.Attempt, nil
	}
	additionalTime, ok := attemptConfig[loginResp.Attempt]
	if !ok {
		// will blocked acccount
		user.Status = userRepo.StatusBlocked
		err = a.userRepo.UpdateStatus(ctx, user.ID, user.Status)
		if err != nil {
			return 0, gers.NewWithMetadata(err,
				pkgLog.Metadata(http.StatusInternalServerError, constant.ErrorMessageTryAgain))
		}
		return 0, gers.NewWithMetadata(errors.New("max attempt login reached"),
			pkgLog.Metadata(http.StatusBadRequest, "Max attempt reached, we have blocked your account"))
	}

	now := time.Now()
	blockTime := now.Add(additionalTime)
	err = a.sessionRepo.SetLoginAttempt(ctx, user.ID, sessionRepo.LoginAttempt{
		Attempt:     loginResp.Attempt + 1,
		BlockTime:   blockTime,
		RefreshTime: blockTime.Add(5 * time.Minute * time.Duration(loginResp.Attempt+1)),
	})
	if err != nil {
		return 0, gers.NewWithMetadata(err,
			pkgLog.Metadata(http.StatusInternalServerError, constant.ErrorMessageTryAgain))
	}

	if blockTime.After(now) {
		return 0, gers.NewWithMetadata(errors.New("too many login attempt"),
			pkgLog.Metadata(http.StatusBadRequest, fmt.Sprintf("Too many login attempt, you account was locked for %s", additionalTime.String())))
	}
	return loginResp.Attempt, nil
}

func (a *Usecase) updateAttemptFailed(ctx context.Context, user userRepo.User, attempt int) (err error) {
	additionalTime, ok := attemptConfig[attempt]
	if !ok {
		// will blocked acccount
		user.Status = userRepo.StatusBlocked
		err = a.userRepo.UpdateStatus(ctx, user.ID, user.Status)
		if err != nil {
			return gers.NewWithMetadata(err,
				pkgLog.Metadata(http.StatusInternalServerError, constant.ErrorMessageTryAgain))
		}
		return gers.NewWithMetadata(errors.New("max attempt login reached"),
			pkgLog.Metadata(http.StatusBadRequest, "Max attempt reached, we have blocked your account"))
	}
	now := time.Now()
	err = a.sessionRepo.SetLoginAttempt(ctx, user.ID, sessionRepo.LoginAttempt{
		Attempt:     attempt + 1,
		BlockTime:   now.Add(additionalTime),
		RefreshTime: now.Add(additionalTime + (5 * time.Minute * time.Duration(attempt))),
	})
	if err != nil {
		return gers.NewWithMetadata(err,
			pkgLog.Metadata(http.StatusInternalServerError, constant.ErrorMessageTryAgain))
	}
	return nil
}

func (a *Usecase) resetAttemptFailed(ctx context.Context, userID int64) error {
	err := a.sessionRepo.SetLoginAttempt(ctx, userID, sessionRepo.LoginAttempt{
		Attempt:     0,
		BlockTime:   time.Now(),
		RefreshTime: time.Now(),
	})
	if err != nil {
		return gers.NewWithMetadata(err,
			pkgLog.Metadata(http.StatusInternalServerError, constant.ErrorMessageTryAgain))
	}
	return nil
}
