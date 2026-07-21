package aws

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/ports"
)

var _ ports.TaskEndpointResolver = (*ECSTaskEndpointResolver)(nil)

type ECSTaskEndpointResolver struct {
	ecsPort ports.ECSPort
	ecsCfg  *configs.ECSConfig

	managementPort int
}

func NewECSTaskEndpointResolver(
	ecsPort ports.ECSPort,
	ecsCfg *configs.ECSConfig,
	managementPort int,
) *ECSTaskEndpointResolver {
	return &ECSTaskEndpointResolver{
		ecsPort:        ecsPort,
		ecsCfg:         ecsCfg,
		managementPort: managementPort,
	}
}

func (r *ECSTaskEndpointResolver) ResolveTaskEndpoint(ctx context.Context, serviceName string, taskID string) (string, error) {
	task, err := r.ecsPort.DescribeTask(
		ctx,
		r.ecsCfg.ClusterName,
		taskID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to describe task: taskID=%s: %w",
			taskID,
			err,
		)
	}

	if task.PrivateIP == "" {
		return "", fmt.Errorf(
			"task private IP is empty: taskID=%s",
			taskID,
		)
	}

	return fmt.Sprintf(
		"%s:%d",
		task.PrivateIP,
		r.managementPort,
	), nil
}
