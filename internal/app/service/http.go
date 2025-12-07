package service

import (
	"context"

	"github.com/aidapedia/jabberwock/internal/interface/http"
	policyRepo "github.com/aidapedia/jabberwock/internal/repository/policy"
	policyUC "github.com/aidapedia/jabberwock/internal/usecase/policy"
)

type ServiceHTTP struct {
	policyUC    policyUC.Interface
	httpService http.HTTPServiceInterface
}

func NewServiceHTTP(httpService http.HTTPServiceInterface, policyUC policyUC.Interface) *ServiceHTTP {
	return &ServiceHTTP{
		policyUC:    policyUC,
		httpService: httpService,
	}
}

func (s *ServiceHTTP) LoadPolicy(ctx context.Context) error {
	err := s.policyUC.LoadPolicy(ctx, policyRepo.HTTPServiceType)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceHTTP) Run() error {
	return s.httpService.ListenAndServe()
}
