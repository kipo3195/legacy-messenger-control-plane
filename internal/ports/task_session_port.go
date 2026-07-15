package ports

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

type TaskSessionPort interface {
	SaveTaskSessionReport(ctx context.Context, report domain.TaskSessionReport) error
	GetTaskSessionReport(ctx context.Context, serviceName string) ([]domain.SessionReport, error)
	GetExpiredReportTask(ctx context.Context, serviceName string) (map[string]string, error)
}
