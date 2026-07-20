package usecase

import (
	"fmt"
	"legacy-messenger-control-plane/internal/domain"
	"sync"
	"time"
)

// Coordinator는 두 스케줄러가 같은 Scale-in 작업 상태를 공유하도록 관리하는 메모리 객체
// SessionAutoScalingScheduler
// - Scale-in 필요 판단
// - Scale-in 작업 등록
// - 현재 Scale-in 진행 여부 조회

// ScaleInScheduler
// - 등록된 작업 조회
// - Drain 진행
// - 작업 상태 변경
// - 완료 또는 실패 처리

// 두 스케줄러가 직접 서로를 호출하지 않고 ScaleInCoordinator를 통해 상태를 공유하는 구조
type ScaleInCoordinator struct {
	mu   sync.RWMutex
	jobs map[string]*domain.ScaleInJob
}

func NewScaleInCoordinator() *ScaleInCoordinator {
	return &ScaleInCoordinator{
		jobs: make(map[string]*domain.ScaleInJob),
	}
}

// 작업 등록
func (c *ScaleInCoordinator) Request(
	job domain.ScaleInJob,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	current, exists := c.jobs[job.ServiceName]
	if exists && isActiveScaleInStatus(current.Status) {
		return fmt.Errorf(
			"scale-in job is already active: serviceName=%s status=%s",
			job.ServiceName,
			current.Status,
		)
	}

	now := time.Now()

	job.Status = domain.ScaleInStatusRequested
	job.RequestedAt = now
	job.UpdatedAt = now

	c.jobs[job.ServiceName] = &job

	return nil
}

// 진행중 작업 조회
func (c *ScaleInCoordinator) GetActive(
	serviceName string,
) (domain.ScaleInJob, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	job, exists := c.jobs[serviceName]
	if !exists || !isActiveScaleInStatus(job.Status) {
		return domain.ScaleInJob{}, false
	}

	return *job, true
}

// 상태 변경
func (c *ScaleInCoordinator) MarkDraining(
	serviceName string,
	targetTaskID string,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	job, exists := c.jobs[serviceName]
	if !exists {
		return fmt.Errorf(
			"scale-in job not found: serviceName=%s",
			serviceName,
		)
	}

	if job.Status != domain.ScaleInStatusRequested {
		return fmt.Errorf(
			"invalid scale-in status transition: current=%s target=%s",
			job.Status,
			domain.ScaleInStatusDraining,
		)
	}

	job.TargetTaskID = targetTaskID
	job.Status = domain.ScaleInStatusDraining
	job.UpdatedAt = time.Now()

	return nil
}

// 세션 0 횟수 증가
func (c *ScaleInCoordinator) IncreaseZeroSessionStreak(
	serviceName string,
) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	job, exists := c.jobs[serviceName]
	if !exists {
		return 0, fmt.Errorf(
			"scale-in job not found: serviceName=%s",
			serviceName,
		)
	}

	job.ZeroSessionStreak++
	job.UpdatedAt = time.Now()

	return job.ZeroSessionStreak, nil
}

// 초기화
func (c *ScaleInCoordinator) ResetZeroSessionStreak(
	serviceName string,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	job, exists := c.jobs[serviceName]
	if !exists {
		return fmt.Errorf(
			"scale-in job not found: serviceName=%s",
			serviceName,
		)
	}

	job.ZeroSessionStreak = 0
	job.UpdatedAt = time.Now()

	return nil
}

// Scale-in 적용 상태 변경
func (c *ScaleInCoordinator) MarkApplied(
	serviceName string,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	job, exists := c.jobs[serviceName]
	if !exists {
		return fmt.Errorf(
			"scale-in job not found: serviceName=%s",
			serviceName,
		)
	}

	if job.Status != domain.ScaleInStatusDraining {
		return fmt.Errorf(
			"scale-in cannot be applied from status=%s",
			job.Status,
		)
	}

	job.Status = domain.ScaleInStatusApplied
	job.UpdatedAt = time.Now()

	return nil
}

// 활성 상태 판단
func isActiveScaleInStatus(
	status domain.ScaleInStatus,
) bool {
	switch status {
	case domain.ScaleInStatusRequested,
		domain.ScaleInStatusDraining,
		domain.ScaleInStatusApplied:
		return true

	default:
		return false
	}
}
