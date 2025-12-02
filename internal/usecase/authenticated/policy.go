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

func (u *Usecase) AddResource(ctx context.Context, req AddResourceRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/AddResource")
	defer span.Finish(err)

	err = u.policyRepo.CreateResource(ctx, policyRepo.Resource{
		Type:   req.Type,
		Method: req.Method,
		Path:   req.Path,
	})
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

func (u *Usecase) AddPermission(ctx context.Context, req AddPermissionRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/AddPermission")
	defer span.Finish(err)

	permission := policyRepo.Permission{
		Name:        req.Name,
		Description: req.Description,
	}
	err = u.policyRepo.CreatePermission(ctx, permission)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if len(req.AssignToResources) > 0 {
		err = u.policyRepo.BulkAssignResources(ctx, permission.ID, req.AssignToResources)
		if err != nil {
			return gers.NewWithMetadata(err,
				ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
		}
	}

	return nil
}

func (u *Usecase) AddRole(ctx context.Context, req AddRoleRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/AddRole")
	defer span.Finish(err)

	role := policyRepo.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	err = u.policyRepo.CreateRole(ctx, role)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if len(req.AssignToPermissions) > 0 {
		err = u.policyRepo.BulkAssignPermissions(ctx, role.ID, req.AssignToPermissions)
		if err != nil {
			return gers.NewWithMetadata(err,
				ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
		}
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
