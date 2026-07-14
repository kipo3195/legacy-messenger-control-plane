package task

import "time"

type TaskSessionValue struct {
	SessionCount int       `json:"sessionCount"`
	ReportedAt   time.Time `json:"reportedAt"`
}
