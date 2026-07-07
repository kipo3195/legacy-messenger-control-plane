package application

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type serviceStatusUsecase struct {
	ecsPort  ports.ECSPort
	ecsCfg   *configs.ECSConfig
	registry *configs.ServiceRegistry
}

type ServiceStatusUsecase interface {
	GetServiceStatus(ctx context.Context, ecsServiceName string) (*domain.ServiceStatus, error)
}

func NewServiceStatusUsecase(ecsPort ports.ECSPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) ServiceStatusUsecase {
	return &serviceStatusUsecase{
		ecsPort:  ecsPort,
		ecsCfg:   ecsCfg,
		registry: registry,
	}
}

func (s *serviceStatusUsecase) GetServiceStatus(ctx context.Context, serviceName string) (*domain.ServiceStatus, error) {

	serviceDef, err := s.registry.Find(serviceName)
	if err != nil {
		return nil, err
	}

	ecsServiceName := serviceDef.ECSServiceName

	return s.ecsPort.DescribeService(ctx, s.ecsCfg.ClusterName, ecsServiceName)
}
