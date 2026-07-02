package bootstrap

import (
	"context"
	"control/legacy-messenger-control-plane/internal/ports"
	"legacy-messenger-control-plane/internal/config"
)

type Clients struct {
	ECS ports.ECSPort
	// CloudWatch ports.MetricPort
	// ELB        ports.TargetHealthPort
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	ecsClient, err := aws.NewECSClient(ctx, cfg.AWS.Region)
	if err != nil {
		return nil, err
	}

	// cloudWatchClient, err := aws.NewCloudWatchClient(ctx, cfg.AWS.Region)
	// if err != nil {
	// 	return nil, err
	// }

	// elbClient, err := aws.NewELBV2Client(ctx, cfg.AWS.Region)
	// if err != nil {
	// 	return nil, err
	// }

	return &Clients{
		ECS: ecsClient,
		// CloudWatch: cloudWatchClient,
		// ELB:        elbClient,
	}, nil
}
