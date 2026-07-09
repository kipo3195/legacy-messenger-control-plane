package bootstrap

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/adapters/aws"
	"legacy-messenger-control-plane/internal/ports"
)

type Clients struct {
	ECS        ports.ECSPort
	CloudWatch ports.CloudWatchPort
	ELB        ports.ELBPort
}

func NewClients(ctx context.Context, cfg *configs.Config) (*Clients, error) {
	ecsClient, err := aws.NewECSClient(ctx, cfg.AWS.Region)
	if err != nil {
		return nil, err
	}

	cloudWatchClient, err := aws.NewCloudWatchClient(ctx, cfg.AWS.Region)
	if err != nil {
		return nil, err
	}

	elbClient, err := aws.NewELBV2Client(ctx, cfg.AWS.Region)
	if err != nil {
		return nil, err
	}

	return &Clients{
		ECS:        ecsClient,
		CloudWatch: cloudWatchClient,
		ELB:        elbClient,
	}, nil
}
