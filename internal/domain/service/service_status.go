package service

type ServiceStatus struct {
    ServiceName    string             `json:"serviceName"`
    ClusterName    string             `json:"clusterName"`
    ECSServiceName string             `json:"ecsServiceName"`
    Status         string             `json:"status"`
    DesiredCount   int32              `json:"desiredCount"`
    RunningCount   int32              `json:"runningCount"`
    PendingCount   int32              `json:"pendingCount"`
    TaskDefinition string             `json:"taskDefinition"`
    Deployments    []DeploymentStatus `json:"deployments"`
    Events         []ServiceEvent     `json:"events"`
}
