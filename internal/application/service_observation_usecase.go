package application

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type serviceObservationUsecase struct {
	ecsPort  ports.ECSPort
	ecsCfg   *configs.ECSConfig
	registry *configs.ServiceRegistry
}

type ServiceObservationUsecase interface {
	GetServiceStatus(ctx context.Context, serviceName string) (*domain.ServiceStatus, error)
	GetServiceList(ctx context.Context) ([]domain.ServiceList, error)
}

func NewServiceObservationUsecase(ecsPort ports.ECSPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) ServiceObservationUsecase {
	return &serviceObservationUsecase{
		ecsPort:  ecsPort,
		ecsCfg:   ecsCfg,
		registry: registry,
	}
}

func (s *serviceObservationUsecase) GetServiceStatus(ctx context.Context, serviceName string) (*domain.ServiceStatus, error) {

	serviceDef, err := s.registry.Find(serviceName)
	if err != nil {
		return nil, err
	}

	ecsServiceName := serviceDef.ECSServiceName
	return s.ecsPort.DescribeService(ctx, s.ecsCfg.ClusterName, ecsServiceName)
}

func (s *serviceObservationUsecase) GetServiceList(ctx context.Context) ([]domain.ServiceList, error) {

	serviceList := make([]domain.ServiceList, 0)
	services := s.registry.List()

	clusterName := s.ecsCfg.ClusterName
	for _, v := range services {

		ecsServiceName := v.ECSServiceName
		serviceStatus, err := s.ecsPort.DescribeService(ctx, clusterName, ecsServiceName)
		if err != nil {

			service := domain.ServiceList{
				ServiceName: ecsServiceName,
				Status:      serviceStatus.Status,
				Deployments: serviceStatus.Deployments,
			}

			serviceList = append(serviceList, service)
		} else {
			return nil, fmt.Errorf("[GetServiceList] %v DescribeService error", ecsServiceName)
		}
	}

	return serviceList, nil
}
