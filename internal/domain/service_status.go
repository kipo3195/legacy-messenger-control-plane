package domain

type ServiceStatus struct {
	ServiceName    string             `json:"serviceName"`
	ECSServiceName string            `json:"ecsServiceName"`
	ClusterName    string            `json:"clusterName"`
	Status         string             `json:"status"`
	DesiredCount  int32              `json:"desiredCount"`
	RunningCount  int32              `json:"runningCount"`
	PendingCount  int32              `json:"pendingCount"`
	TaskDefinition string             `json:"taskDefinition"`
	Deployments   []DeploymentStatus `json:"deployments"`
	Events        []ServiceEvent     `json:"events"`
}


