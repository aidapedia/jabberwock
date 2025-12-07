package policy

import (
	"context"
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
		u.enforcer.AddPolicy(StdPolicy(policy)...)
	}

	return nil
}

func (u *Usecase) GetUserPermissions(ctx context.Context, userID int64) (resp GetUserPermissionsResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/GetUserPermissions")
	defer span.Finish(err)

	role, err := u.policyRepo.GetRoleByUserID(ctx, userID)
	if err != nil {
		return GetUserPermissionsResponse{}, gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	if role == nil {
		return GetUserPermissionsResponse{}, gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	var permissions []policyRepo.Permission
	if len(role) > 0 && role[0].ID == policyRepo.SuperAdminRole {
		permissions, err = u.policyRepo.GetAllPermissions(ctx)
		if err != nil {
			return GetUserPermissionsResponse{}, gers.NewWithMetadata(err,
				ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
		}
	} else {
		permissions, err = u.policyRepo.GetUserPermissions(ctx, userID)
		if err != nil {
			return GetUserPermissionsResponse{}, gers.NewWithMetadata(err,
				ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
		}
	}

	resp.Transform(permissions)
	return resp, nil
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
		err = u.policyRepo.BulkAssignResources(ctx, nil, permission.ID, req.AssignToResources)
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
		err = u.policyRepo.BulkAssignPermissions(ctx, nil, role.ID, req.AssignToPermissions)
		if err != nil {
			return gers.NewWithMetadata(err,
				ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
		}
	}

	return nil
}

func (u *Usecase) UpdateResource(ctx context.Context, req UpdateResourceRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/UpdateResource")
	defer span.Finish(err)

	err = u.policyRepo.UpdateResource(ctx, policyRepo.Resource{
		ID:     req.ID,
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

func (u *Usecase) UpdatePermission(ctx context.Context, req UpdatePermissionRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/UpdatePermission")
	defer span.Finish(err)

	err = u.policyRepo.UpdatePermission(ctx, policyRepo.Permission{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

func (u *Usecase) UpdateRole(ctx context.Context, req UpdateRoleRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/UpdateRole")
	defer span.Finish(err)

	err = u.policyRepo.UpdateRole(ctx, policyRepo.Role{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

func (u *Usecase) DeleteResource(ctx context.Context, req DeleteResourceRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/DeleteResource")
	defer span.Finish(err)

	err = u.policyRepo.DeleteResource(ctx, req.ID)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

func (u *Usecase) DeletePermission(ctx context.Context, req DeletePermissionRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/DeletePermission")
	defer span.Finish(err)

	err = u.policyRepo.DeletePermission(ctx, req.ID)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}

func (u *Usecase) DeleteRole(ctx context.Context, req DeleteRoleRequest) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "AuthenticateUsecase/DeleteRole")
	defer span.Finish(err)

	err = u.policyRepo.DeleteRole(ctx, req.ID)
	if err != nil {
		return gers.NewWithMetadata(err,
			ghttp.Metadata(http.StatusInternalServerError, cerror.ErrorMessageTryAgain.Error()))
	}

	return nil
}
