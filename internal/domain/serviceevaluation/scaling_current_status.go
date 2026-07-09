package serviceevaluation

type ScalingCurrentStatus struct {
	ActiveConnectionCount float64 `json:"activeConnectionCount"`
	RunningTaskCount      int     `json:"runningTaskCount"`
	DesiredCount          int     `json:"desiredCount"`
	ConnectionPerTask     float64 `json:"connectionPerTask"`
}
