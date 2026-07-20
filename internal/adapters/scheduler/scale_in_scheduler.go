package scheduler

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/internal/application/usecase"
	"time"
)

type ScaleInScheduler struct {
	usecase  usecase.ScaleInUsecase
	interval time.Duration
}

func NewScaleInScheduler(usecase usecase.ScaleInUsecase, interval time.Duration) *ScaleInScheduler {
	return &ScaleInScheduler{
		usecase:  usecase,
		interval: interval,
	}
}

func (s *ScaleInScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := s.usecase.Process(ctx); err != nil {
				fmt.Printf(
					"[scale-in scheduler] process failed: %v\n",
					err,
				)
			}
		}
	}
}
