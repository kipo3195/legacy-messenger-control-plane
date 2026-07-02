package application

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

type serviceStatusUsecase struct {
	ecsPort ECSPort
}

type ServiceStatusUsecase interface {
	GetServiceStatus(clusterName string, ecsServiceName string) (*domain.ServiceStatus, error)
}

func NewServiceStatusUsecase(ecsPort ECSPort) ServiceStatusUsecase {
	return &serviceStatusUsecase{
		ecsPort: ecsPort,
	}
}

func (s *serviceStatusUsecase) GetServiceStatus(clusterName string, ecsServiceName string) (*domain.ServiceStatus, error) {
	return s.ecsPort.DescribeService(context.Background(), clusterName, ecsServiceName)
}
