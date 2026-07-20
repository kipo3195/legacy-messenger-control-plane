package scheduler

import (
	"context"
	"fmt"
	"time"
)

type ScaleInScheduler struct {
	usecase  ScaleInUsecase
	interval time.Duration
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
