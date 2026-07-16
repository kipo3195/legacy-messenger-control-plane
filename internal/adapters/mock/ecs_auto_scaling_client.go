package mock

import (
	"context"
	"sync"

	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type SessionAutoScalingECSClient struct {
	mu sync.Mutex

	State domain.ECSServiceControlState

	GetStateError error
	UpdateError   error

	UpdateCalled        bool
	UpdatedDesiredCount int
}

var _ ports.SessionAutoScalingECSPort = (*SessionAutoScalingECSClient)(nil)

func NewSessionAutoScalingECSClient(
	initialState domain.ECSServiceControlState,
) *SessionAutoScalingECSClient {
	return &SessionAutoScalingECSClient{
		State: initialState,
	}
}

func (c *SessionAutoScalingECSClient) GetServiceControlState(
	ctx context.Context,
	clusterName string,
	ecsServiceName string,
) (domain.ECSServiceControlState, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.GetStateError != nil {
		return domain.ECSServiceControlState{}, c.GetStateError
	}

	return c.State, nil
}

func (c *SessionAutoScalingECSClient) UpdateServiceDesiredCount(
	ctx context.Context,
	clusterName string,
	ecsServiceName string,
	desiredCount int,
) (domain.ECSServiceControlState, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.UpdateCalled = true
	c.UpdatedDesiredCount = desiredCount

	if c.UpdateError != nil {
		return domain.ECSServiceControlState{}, c.UpdateError
	}

	c.State.DesiredCount = int32(desiredCount)

	return c.State, nil
}
