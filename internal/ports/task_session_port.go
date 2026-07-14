package ports

import (
	"context"
	"legacy-messenger-control-plane/internal/domain"
)

type TaskSessionPort interface {
	SaveTaskSessionReport(ctx context.Context, report domain.TaskSessionReport) error
}
