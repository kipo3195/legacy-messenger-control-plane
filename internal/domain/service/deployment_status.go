package service

type DeploymentStatus struct {
    Status       string `json:"status"`
    DesiredCount int32  `json:"desiredCount"`
    RunningCount int32  `json:"runningCount"`
    PendingCount int32  `json:"pendingCount"`
    RolloutState string `json:"rolloutState"`
}
