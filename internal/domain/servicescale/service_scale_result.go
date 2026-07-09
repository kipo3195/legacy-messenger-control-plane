package servicescale

type ServiceScaleResult struct {
	ServiceName          string `json:"serviceName"`
	ECSServiceName       string `json:"ecsServiceName"`
	ClusterName          string `json:"clusterName"`
	PreviousDesiredCount int32  `json:"previousDesiredCount"`
	DesiredCount         int32  `json:"desiredCount"`
	RunningCount         int32  `json:"runningCount"`
	PendingCount         int32  `json:"pendingCount"`
	Status               string `json:"status"`
	Message              string `json:"message"`
}
