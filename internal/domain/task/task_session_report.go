package task

import "time"

type TaskSessionReport struct {
	ServiceName  string
	TaskID       string
	SessionCount int
	ReportedAt   time.Time
	ExpiresAt    time.Time
}
