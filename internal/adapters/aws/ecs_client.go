package aws

import (
	"context"
	"fmt"
	"time"

	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var _ ports.ECSPort = (*ECSClient)(nil)

type ECSClient struct {
	client *ecs.Client
}

func NewECSClient(ctx context.Context, region string) (*ECSClient, error) {
	if region == "" {
		return nil, fmt.Errorf("aws region is required")
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &ECSClient{
		client: ecs.NewFromConfig(cfg),
	}, nil
}

func (c *ECSClient) DescribeService(ctx context.Context, clusterName string, ecsServiceName string) (*domain.ServiceStatus, error) {

	if clusterName == "" {
		return nil, fmt.Errorf("clusterName is required")
	}

	if ecsServiceName == "" {
		return nil, fmt.Errorf("ecsServiceName is required")
	}

	out, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: []string{ecsServiceName},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe ecs service %s: %w", ecsServiceName, err)
	}

	if len(out.Services) == 0 {
		return nil, fmt.Errorf("ecs service not found: %s", ecsServiceName)
	}

	svc := out.Services[0]

	status := &domain.ServiceStatus{
		ServiceName:    ecsServiceName,
		ClusterName:    clusterName,
		Status:         ptrString(svc.Status),
		DesiredCount:   svc.DesiredCount,
		RunningCount:   svc.RunningCount,
		PendingCount:   svc.PendingCount,
		TaskDefinition: ptrString(svc.TaskDefinition),
	}

	for _, d := range svc.Deployments {
		status.Deployments = append(status.Deployments, domain.DeploymentStatus{
			Status:       ptrString(d.Status),
			DesiredCount: d.DesiredCount,
			RunningCount: d.RunningCount,
			PendingCount: d.PendingCount,
			RolloutState: string(d.RolloutState),
		})
	}

	for _, e := range svc.Events {
		status.Events = append(status.Events, domain.ServiceEvent{
			Message:   ptrString(e.Message),
			CreatedAt: ptrTime(e.CreatedAt),
		})
	}

	return status, nil
}

func ptrString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func ptrTime(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}
