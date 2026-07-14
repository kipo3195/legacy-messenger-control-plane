package usecase

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type sessionAutoScalingUsecase struct {
	taskSessionPort ports.TaskSessionPort
	ecsPort         ports.ECSPort
	registry        *configs.ServiceRegistry
}

type SessionAutoScalingUsecase interface {
	EvaluateAndScale(ctx context.Context, serviceName string) (domain.SessionAutoScalingResult, error)
}

func NewSessionAutoScalingUsecase(
	taskSessionPort ports.TaskSessionPort,
	ecsPort ports.ECSPort,
	registry *configs.ServiceRegistry,
	scalingConfig *configs.ScalingConfig,
) SessionAutoScalingUsecase {
	return &sessionAutoScalingUsecase{
		taskSessionPort: taskSessionPort,
		ecsPort:         ecsPort,
		registry:        registry,
	}
}

func (u *sessionAutoScalingUsecase) EvaluateAndScale(ctx context.Context, serviceName string) (domain.SessionAutoScalingResult, error) {

	// redis 에서 값 조회

	u.taskSessionPort.GetTaskSessionReport(ctx, serviceName)

	// task 값을 기준으로 scaling 판단

	// desired count 변경 처리

	return domain.SessionAutoScalingResult{}, nil
}
