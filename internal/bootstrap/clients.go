package bootstrap

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/adapters/aws"
	"legacy-messenger-control-plane/internal/adapters/mock"
	"legacy-messenger-control-plane/internal/adapters/redis"
	"legacy-messenger-control-plane/internal/adapters/ssh"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type Clients struct {
	ECS            ports.ECSPort
	AutoScalingECS ports.SessionAutoScalingECSPort // mock
	CloudWatch     ports.CloudWatchPort
	ELB            ports.ELBPort
	TaskSession    ports.TaskSessionPort

	closeRedis func() error
}

func NewClients(ctx context.Context, cfg *configs.Config) (*Clients, error) {
	ecsClient, err := aws.NewECSClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	mockECSClient := mock.NewSessionAutoScalingECSClient(
		domain.ECSServiceControlState{
			DesiredCount: 5,
			RunningCount: 2,
			PendingCount: 0,
		},
	)

	cloudWatchClient, err := aws.NewCloudWatchClient(ctx, cfg.AWS.Region)
	if err != nil {
		return nil, err
	}

	elbClient, err := aws.NewELBV2Client(ctx, cfg.AWS.Region)
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.NewSSHClient(cfg.SSH)

	taskSessionClient, err := redis.NewRedisClient(ctx, cfg.Redis, sshClient)
	if err != nil {
		return nil, err
	}

	return &Clients{
		ECS:            ecsClient,
		AutoScalingECS: mockECSClient,
		CloudWatch:     cloudWatchClient,
		ELB:            elbClient,
		TaskSession:    taskSessionClient,
		closeRedis:     taskSessionClient.Close,
	}, nil
}

// 외부에서 서비스 종료시 호출
func (c *Clients) Close() error {
	if c.closeRedis != nil {
		return c.closeRedis()
	}

	return nil
}
