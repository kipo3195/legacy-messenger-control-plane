package client

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/internal/ports"
	"net/http"
	"time"
)

var _ ports.TaskDrainPort = (*TaskDrainClient)(nil)

type TaskDrainClient struct {
	httpClient *http.Client
	resolver   ports.TaskEndpointResolver
}

func NewTaskDrainClient(
	resolver ports.TaskEndpointResolver,
) *TaskDrainClient {
	return &TaskDrainClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		resolver: resolver,
	}
}
func (c *TaskDrainClient) RequestDrain(
	ctx context.Context,
	serviceName string,
	taskID string,
) error {

	endpoint, err := c.resolver.ResolveTaskEndpoint(
		ctx,
		serviceName,
		taskID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to resolve task endpoint: serviceName=%s taskID=%s: %w",
			serviceName,
			taskID,
			err,
		)
	}

	requestURL := fmt.Sprintf(
		"http://%s/internal/v1/drain",
		endpoint,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		requestURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create task drain request: taskID=%s: %w",
			taskID,
			err,
		)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"failed to call task drain API: taskID=%s endpoint=%s: %w",
			taskID,
			endpoint,
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return fmt.Errorf(
			"task drain API returned unexpected status: taskID=%s statusCode=%d",
			taskID,
			resp.StatusCode,
		)
	}

	return nil
}

func (c *TaskDrainClient) CancelDrain(
	ctx context.Context,
	serviceName string,
	taskID string,
) error {

	endpoint, err := c.resolver.ResolveTaskEndpoint(
		ctx,
		serviceName,
		taskID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to resolve task endpoint: serviceName=%s taskID=%s: %w",
			serviceName,
			taskID,
			err,
		)
	}

	requestURL := fmt.Sprintf(
		"http://%s/internal/v1/drain/cancel",
		endpoint,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		requestURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create drain cancel request: taskID=%s: %w",
			taskID,
			err,
		)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"failed to call drain cancel API: taskID=%s: %w",
			taskID,
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return fmt.Errorf(
			"drain cancel API returned unexpected status: taskID=%s statusCode=%d",
			taskID,
			resp.StatusCode,
		)
	}

	return nil
}
