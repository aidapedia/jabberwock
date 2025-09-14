package service

import (
	"context"

	"github.com/kurniajigunawan/homestay/internal/interface/http"
	policyRepo "github.com/kurniajigunawan/homestay/internal/repository/policy"
	policyUC "github.com/kurniajigunawan/homestay/internal/usecase/policy"
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
