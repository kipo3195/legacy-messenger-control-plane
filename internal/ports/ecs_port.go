package ports

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

type ECSPort interface {
	DescribeService(ctx context.Context, clusterName string, ecsServiceName string) (*domain.ServiceStatus, error)
	DescribeTask(ctx context.Context, clusterName string, ecsServiceName string, desiredStatus string) ([]domain.TaskStatus, error)
	GetServiceTargetGroups(ctx context.Context, clusterName string, ecsServiceName string) ([]domain.ServiceTargetGroup, error)
}
