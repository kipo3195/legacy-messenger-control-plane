package service

type ServiceList struct {
	ServiceName    string             `json:"serviceName"`
	ECSServiceName string             `json:"ecsServiceName"`
	Status         string             `json:"status"`
	Deployments    []DeploymentStatus `json:"deployments"`
}
