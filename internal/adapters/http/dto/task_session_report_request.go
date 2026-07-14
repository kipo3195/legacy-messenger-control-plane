package dto

import "legacy-messenger-control-plane/internal/application/command"

type TaskSessionReportRequest struct {
	SessionCount *int `json:"sessionCount" binding:"required,gte=0"`
}

func (r TaskSessionReportRequest) ToCommand(serviceName string, taskID string, sessionCount int) command.TaskSessionReportCommand {
	return command.TaskSessionReportCommand{
		ServiceName:  serviceName,
		TaskID:       taskID,
		SessionCount: sessionCount,
	}
}
