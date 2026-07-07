package domain

type ServiceList struct {
	ServiceName string             `json:"serviceName"`
	Status      string             `json:"status"`
	Deployments []DeploymentStatus `json:"deployments"`
}
