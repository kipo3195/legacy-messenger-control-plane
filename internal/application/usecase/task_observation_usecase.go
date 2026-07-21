package usecase

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type taskObservationUsecase struct {
	ecsPort  ports.ECSPort
	ecsCfg   *configs.ECSConfig
	registry *configs.ServiceRegistry
}

type TaskObservationUsecase interface {
	GetTaskStatus(ctx context.Context, serviceName string, desiredStatus string) ([]domain.TaskStatus, error)
}

func NewTaskObservationUsecase(ecsPort ports.ECSPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) TaskObservationUsecase {
	return &taskObservationUsecase{
		ecsPort:  ecsPort,
		ecsCfg:   ecsCfg,
		registry: registry,
	}
}

func (s *taskObservationUsecase) GetTaskStatus(ctx context.Context, serviceName string, desiredStatus string) ([]domain.TaskStatus, error) {

	serviceDef, err := s.registry.Find(serviceName)
	if err != nil {
		return nil, err
	}

	ecsServiceName := serviceDef.ECSServiceName
	if ecsServiceName == "" {
		return nil, fmt.Errorf("ecsServiceName is empty for service: %s", serviceName)
	}

	taskStatus, err := s.ecsPort.DescribeTasks(ctx, s.ecsCfg.ClusterName, ecsServiceName, desiredStatus)
	if err != nil {
		return nil, err
	}

	// taskStatus 파싱

	return taskStatus, nil
}
