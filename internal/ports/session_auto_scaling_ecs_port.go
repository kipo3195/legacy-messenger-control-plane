package ports

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

// auto scaling 테스트용 Mock생성을 위한 별도 port

type SessionAutoScalingECSPort interface {
	GetServiceControlState(
		ctx context.Context,
		clusterName string,
		ecsServiceName string,
	) (domain.ECSServiceControlState, error)

	UpdateServiceDesiredCount(
		ctx context.Context,
		clusterName string,
		ecsServiceName string,
		desiredCount int,
	) (domain.ECSServiceControlState, error)
}
