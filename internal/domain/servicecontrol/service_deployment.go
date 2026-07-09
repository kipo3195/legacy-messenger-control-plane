package servicecontrol

import "time"

type ServiceDeployment struct {
	ID                 string     `json:"id"`
	Status             string     `json:"status"`
	RolloutState       string     `json:"rolloutState,omitempty"`
	RolloutStateReason string     `json:"rolloutStateReason,omitempty"`
	TaskDefinition     string     `json:"taskDefinition,omitempty"`
	DesiredCount       int32      `json:"desiredCount"`
	RunningCount       int32      `json:"runningCount"`
	PendingCount       int32      `json:"pendingCount"`
	CreatedAt          *time.Time `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
}
