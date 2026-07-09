package domain

type ConnectionPressure struct {
	ServiceName    string `json:"serviceName"`
	ECSServiceName string `json:"ecsServiceName"`
	ClusterName    string `json:"clusterName"`

	ActiveConnectionCount float64 `json:"activeConnectionCount"`
	RunningTaskCount      int     `json:"runningTaskCount"`
	DesiredCount          int     `json:"desiredCount"`

	ConnectionPerTask        float64 `json:"connectionPerTask"`
	TargetConnectionsPerTask int     `json:"targetConnectionsPerTask"`

	PressureStatus string `json:"pressureStatus"`

	ScalingRecommendation ScalingRecommendation `json:"scalingRecommendation"`

	Metric ConnectionPressureMetric `json:"metric"`
}
