package serviceevaluation

type ScalingPolicyStatus struct {
	TargetConnectionsPerTask int `json:"targetConnectionsPerTask"`
	ScaleOutThreshold        int `json:"scaleOutThreshold"`
	ScaleInThreshold         int `json:"scaleInThreshold"`
	MinCount                 int `json:"minCount"`
	MaxCount                 int `json:"maxCount"`
}
