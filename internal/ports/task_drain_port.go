package ports

import "context"

type TaskDrainPort interface {
	RequestDrain(ctx context.Context, serviceName string, taskID string) error
	CancelDrain(ctx context.Context, serviceName string, taskID string) error
}
