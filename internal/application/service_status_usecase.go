package application

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type serviceStatusUsecase struct {
	ecsPort ports.ECSPort
}

type ServiceStatusUsecase interface {
	GetServiceStatus(clusterName string, ecsServiceName string) (*domain.ServiceStatus, error)
}

func NewServiceStatusUsecase(ecsPort ports.ECSPort) ServiceStatusUsecase {
	return &serviceStatusUsecase{
		ecsPort: ecsPort,
	}
}

func (s *serviceStatusUsecase) GetServiceStatus(clusterName string, ecsServiceName string) (*domain.ServiceStatus, error) {
	return s.ecsPort.DescribeService(context.Background(), clusterName, ecsServiceName)
}
