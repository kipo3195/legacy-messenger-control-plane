package bootstrap

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/adapters/scheduler"
	"time"
)

type Scheduler struct {
	AutoScale scheduler.SessionScalingScheduler
}

func NewScheduler(ctx context.Context, usecases *UseCases, cfg *configs.Config) error {

	interval := cfg.AutoScale.Interval
	// time.Duration은 내부적으로 int64 기반이지만, Go는 int와 time.Duration 사이를 자동 변환하지 않는다.
	autoScale := scheduler.NewSessionScalingScheduler(usecases.AutoScale, "ws", time.Duration(interval)*time.Second)
	go autoScale.Start(ctx)

	scaleIn := scheduler.NewScaleInScheduler(usecases.ScaleInUsecase, time.Duration(interval)*time.Second)
	go scaleIn.Start(ctx)
	return nil
}
