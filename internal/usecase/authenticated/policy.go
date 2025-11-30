package authenticated

import (
	"context"
	"fmt"
	"net/http"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/telemetry/tracer"
	cerror "github.com/aidapedia/jabberwock/internal/common/error"
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
)

func (u *Usecase) LoadPolicy(ctx context.Context, serviceType policyRepo.ServiceType) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/LoadPolicy")
	defer span.Finish(err)

	policies, err := u.policyRepo.LoadPolicy(ctx, serviceType)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	u.enforcer.ClearPolicy()

	for _, policy := range policies {
		u.enforcer.AddPolicy(stdPolicy(policy)...)
	}

	return nil
}

// Will Generated Policy Like
// 1, http:GET, /api/v1/users
// 1, rpc:OrderService, GetOrder
func stdPolicy(policy policyRepo.Policy) []interface{} {
	return []interface{}{
		policy.Role,
		fmt.Sprintf("%s:%s", policy.Type, policy.Method),
		policy.Path,
	}
}
