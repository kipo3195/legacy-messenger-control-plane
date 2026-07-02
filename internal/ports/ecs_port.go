package ports

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

type ECSPort interface {
	DescribeService(ctx context.Context, clusterName string, ecsServiceName string) (*domain.ServiceStatus, error)
}
