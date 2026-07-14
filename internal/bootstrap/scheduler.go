package bootstrap

import (
	"context"
	"legacy-messenger-control-plane/internal/adapters/scheduler"
)

type Scheduler struct {
	AutoScale scheduler.SessionScalingScheduler
}

func NewScheduler(ctx context.Context, usecases *UseCases) error {

	autoScale := scheduler.NewSessionScalingScheduler(usecases.AutoScale, "ws", 30)
	go autoScale.Start(ctx)
	return nil
}
