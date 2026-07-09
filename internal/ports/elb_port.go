package ports

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

type ELBPort interface {
	DescribeTargetHealth(ctx context.Context, targetGroupArn string, loadBalancerType string) (*domain.TargetGroupHealth, error)
	GetLoadBalancerArnByTargetGroupArn(ctx context.Context, targetGroupArn string) (string, error)
}
