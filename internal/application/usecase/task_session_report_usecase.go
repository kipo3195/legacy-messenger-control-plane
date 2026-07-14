package usecase

import (
	"context"
	"legacy-messenger-control-plane/internal/application/command"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
	"time"
)

type taskSessionReportUsecase struct {
	taskSessionPort ports.TaskSessionPort
}

type TaskSessionReportUsecase interface {
	PutTaskSessionReport(ctx context.Context, cmd command.TaskSessionReportCommand) (domain.TaskSessionReportResult, error)
}

func NewTaskSessionReportUsecase(taskSessionPort ports.TaskSessionPort) TaskSessionReportUsecase {
	return &taskSessionReportUsecase{
		taskSessionPort: taskSessionPort,
	}
}

func (h *taskSessionReportUsecase) PutTaskSessionReport(ctx context.Context, cmd command.TaskSessionReportCommand) (domain.TaskSessionReportResult, error) {

	reportedAt := time.Now()

	// 만료시간 서버 설정화
	expiresAt := reportedAt.Add(60 * time.Second)
	report := domain.TaskSessionReport{
		ServiceName:  cmd.ServiceName,
		TaskID:       cmd.TaskID,
		SessionCount: cmd.SessionCount,
		ReportedAt:   reportedAt,
		ExpiresAt:    expiresAt,
	}

	err := h.taskSessionPort.SaveTaskSessionReport(ctx, report)
	if err != nil {
		return domain.TaskSessionReportResult{}, err
	}

	return domain.TaskSessionReportResult{
		ServiceName:  report.ServiceName,
		TaskID:       report.TaskID,
		SessionCount: report.SessionCount,
		ReportedAt:   report.ReportedAt,
		ExpiresAt:    report.ExpiresAt,
	}, nil
}
