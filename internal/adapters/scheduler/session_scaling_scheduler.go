package scheduler

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/internal/application/usecase"
	"log"
	"time"
)

// handler와 동일하게 외부 -> 내부 usecase 진입점

type SessionScalingScheduler struct {
	usecase     usecase.SessionAutoScalingUsecase
	serviceName string
	interval    time.Duration
}

func NewSessionScalingScheduler(
	usecase usecase.SessionAutoScalingUsecase,
	serviceName string,
	interval time.Duration,
) *SessionScalingScheduler {
	return &SessionScalingScheduler{
		usecase:     usecase,
		serviceName: serviceName,
		interval:    interval,
	}
}

func (s *SessionScalingScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// 서비스 시작 직후 한 번 실행하고 싶을 때
	s.execute(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.execute(ctx)
		}
	}
}

func (s *SessionScalingScheduler) execute(ctx context.Context) {
	result, err := s.usecase.EvaluateAndScale(ctx, s.serviceName)
	if err != nil {
		log.Printf(
			"failed to evaluate session scaling: service=%s error=%v",
			s.serviceName,
			err,
		)
	}
	// 이력관리 혹은 로깅
	fmt.Println(result)
}
