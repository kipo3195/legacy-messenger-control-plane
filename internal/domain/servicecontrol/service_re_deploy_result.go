package servicecontrol

type ServiceRedeployResult struct {
	ServiceName    string              `json:"serviceName"`
	ECSServiceName string              `json:"ecsServiceName"`
	ClusterName    string              `json:"clusterName"`
	Action         string              `json:"action"`
	Status         string              `json:"status"`
	Reason         string              `json:"reason,omitempty"`
	DesiredCount   int32               `json:"desiredCount"`
	RunningCount   int32               `json:"runningCount"`
	PendingCount   int32               `json:"pendingCount"`
	Deployments    []ServiceDeployment `json:"deployments"`
	Message        string              `json:"message"`
}
