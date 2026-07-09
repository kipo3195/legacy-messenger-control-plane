package targethealth

type TargetHealthOverallStatus string

const (
    TargetHealthOverallHealthy       TargetHealthOverallStatus = "HEALTHY"
    TargetHealthOverallDegraded      TargetHealthOverallStatus = "DEGRADED"
    TargetHealthOverallTransitioning TargetHealthOverallStatus = "TRANSITIONING"
    TargetHealthOverallUnknown       TargetHealthOverallStatus = "UNKNOWN"
)

type TargetHealthResponse struct {
    ServiceName    string                    `json:"serviceName"`
    ECSServiceName string                    `json:"ecsServiceName"`
    ClusterName    string                    `json:"clusterName"`
    OverallStatus  TargetHealthOverallStatus `json:"overallStatus"`
    Summary        TargetHealthSummary       `json:"summary"`
    TargetGroups   []TargetGroupHealth       `json:"targetGroups"`
}
