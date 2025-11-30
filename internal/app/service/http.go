package service

import (
	"context"

	"github.com/aidapedia/jabberwock/internal/interface/http"
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	authenticatedUC "github.com/aidapedia/jabberwock/internal/usecase/authenticated"
)

type ServiceHTTP struct {
	authenticatedUC authenticatedUC.Interface
	httpService     http.HTTPServiceInterface
}

func NewServiceHTTP(httpService http.HTTPServiceInterface, authenticatedUC authenticatedUC.Interface) *ServiceHTTP {
	return &ServiceHTTP{
		authenticatedUC: authenticatedUC,
		httpService:     httpService,
	}
}

func (s *ServiceHTTP) LoadPolicy(ctx context.Context) error {
	err := s.authenticatedUC.LoadPolicy(ctx, policyRepo.HTTPServiceType)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceHTTP) Run() error {
	return s.httpService.ListenAndServe()
}
