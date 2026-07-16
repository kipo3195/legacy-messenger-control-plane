package fake

import (
	"context"
	"fmt"
	"sync"

	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type ECSClient struct {
	mu sync.RWMutex

	serviceStates map[string]domain.ECSServiceControlState
	redeployCount map[string]int
}

var _ ports.ECSPort = (*ECSClient)(nil)

func NewECSClient(
	initialStates map[string]domain.ECSServiceControlState,
) *ECSClient {
	states := make(map[string]domain.ECSServiceControlState)

	for serviceName, state := range initialStates {
		states[serviceName] = state
	}

	return &ECSClient{
		serviceStates: states,
		redeployCount: make(map[string]int),
	}
}

func (c *ECSClient) GetServiceControlState(
	ctx context.Context,
	clusterName string,
	serviceName string,
) (domain.ECSServiceControlState, error) {
	if err := ctx.Err(); err != nil {
		return domain.ECSServiceControlState{}, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	state, exists := c.serviceStates[serviceName]
	if !exists {
		return domain.ECSServiceControlState{},
			fmt.Errorf("fake ECS service not found: %s", serviceName)
	}

	return state, nil
}

func (c *ECSClient) DescribeService(ctx context.Context, clusterName string, ecsServiceName string) (*domain.ServiceStatus, error) {
	return nil, nil
}
func (c *ECSClient) DescribeTask(ctx context.Context, clusterName string, ecsServiceName string, desiredStatus string) ([]domain.TaskStatus, error) {
	return nil, nil
}

func (c *ECSClient) GetServiceTargetGroups(ctx context.Context, clusterName string, ecsServiceName string) ([]domain.ServiceTargetGroup, error) {
	return nil, nil
}

func (c *ECSClient) GetServiceTargetGroupArn(ctx context.Context, clusterName string, ecsServiceName string) (string, error) {
	return "", nil
}

func (c *ECSClient) UpdateServiceDesiredCount(
	ctx context.Context,
	clusterName string,
	serviceName string,
	desiredCount int,
) (domain.ECSServiceControlState, error) {
	if err := ctx.Err(); err != nil {
		return domain.ECSServiceControlState{}, err
	}

	if desiredCount < 0 {
		return domain.ECSServiceControlState{},
			fmt.Errorf("desired count must not be negative: %d", desiredCount)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	state, exists := c.serviceStates[serviceName]
	if !exists {
		return domain.ECSServiceControlState{},
			fmt.Errorf("fake ECS service not found: %s", serviceName)
	}

	state.DesiredCount = int32(desiredCount)

	// Fake 환경에서는 스케일링이 즉시 완료된 것으로 처리합니다.
	state.RunningCount = int32(desiredCount)
	state.PendingCount = 0

	c.serviceStates[serviceName] = state

	return state, nil
}
func (c *ECSClient) ForceNewDeployment(ctx context.Context, clusterName string, ecsServiceName string) (domain.ServiceRedeployResult, error) {
	return domain.ServiceRedeployResult{}, nil
}
