package task

import "time"

type TaskSessionReportResult struct {
	ServiceName  string    `json:"serviceName"`
	TaskId       string    `json:"taskId"`
	SessionCount int       `json:"sessionCount"`
	ReportedAt   time.Time `json:"reportedAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
