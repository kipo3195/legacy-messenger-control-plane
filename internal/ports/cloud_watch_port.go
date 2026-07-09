package ports

import (
	"context"
	"time"
)

type CloudWatchPort interface {
	GetALBActiveConnectionCount(ctx context.Context, loadBalancerArn string, periodSeconds int32, lookback time.Duration) (float64, error)
}
