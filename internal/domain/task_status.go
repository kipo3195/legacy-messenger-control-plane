package domain

import "time"

type TaskStatus struct {
	TaskID               string `json:"taskId"`
	TaskArn              string `json:"taskArn"`
	TaskDefinition       string `json:"taskDefinition"`
	LastStatus           string `json:"lastStatus"`
	DesiredStatus        string `json:"desiredStatus"`
	HealthStatus         string `json:"healthStatus"`
	LaunchType           string `json:"launchType"`
	AvailabilityZone     string `json:"availabilityZone"`
	ContainerInstanceArn string `json:"containerInstanceArn,omitempty"`
	CapacityProviderName string `json:"capacityProviderName,omitempty"`

	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	StoppingAt *time.Time `json:"stoppingAt,omitempty"`
	StoppedAt  *time.Time `json:"stoppedAt,omitempty"`

	StopCode      string `json:"stopCode,omitempty"`
	StoppedReason string `json:"stoppedReason,omitempty"`

	Containers []ContainerStatus `json:"containers"`
	Network    TaskNetworkInfo   `json:"network"`
}
