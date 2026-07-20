package ports

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"time"
)

type TaskSessionPort interface {
	SaveTaskSessionReport(ctx context.Context, report domain.TaskSessionReport) error
	GetTaskSessionReport(ctx context.Context, serviceName string) (map[string]domain.SessionReport, error)
	GetInvalidReportTask(ctx context.Context, serviceName string, cfg *configs.AutoScaleConfig) (map[string]string, []string, error)
	ShouldStopTask(ctx context.Context, serviceName string, taskID string, now time.Time) (bool, error)
	DeleteTaskSessionState(ctx context.Context, serviceName string, taskID string) error
	GetTaskSessionReportByTask(ctx context.Context, serviceName string, taskID string) (domain.SessionReport, error)
}
