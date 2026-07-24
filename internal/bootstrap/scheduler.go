package bootstrap

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/adapters/scheduler"
	"legacy-messenger-control-plane/logger"
	"time"
)

type Scheduler struct {
	AutoScale scheduler.SessionScalingScheduler
}

func NewScheduler(ctx context.Context, usecases *UseCases, cfg *configs.Config) error {

	interval := cfg.AutoScale.Interval

	seoulLocation, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return fmt.Errorf("failed to load Asia/Seoul timezone: %w", err)
	}
	scalingResultLogger := logger.NewScalingResultLogger(
		"log",
		seoulLocation,
	)

	// time.Duration은 내부적으로 int64 기반이지만, Go는 int와 time.Duration 사이를 자동 변환하지 않는다.
	autoScale := scheduler.NewSessionScalingScheduler(usecases.AutoScale, "ws", time.Duration(interval)*time.Second, scalingResultLogger)
	go autoScale.Start(ctx)

	scaleIn := scheduler.NewScaleInScheduler(usecases.ScaleInUsecase, time.Duration(interval)*time.Second)
	go scaleIn.Start(ctx)
	return nil
}
