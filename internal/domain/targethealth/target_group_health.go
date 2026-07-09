package targethealth

type TargetGroupHealth struct {
    TargetGroupArn   string              `json:"targetGroupArn"`
    TargetGroupName  string              `json:"targetGroupName,omitempty"`
    LoadBalancerType string              `json:"loadBalancerType,omitempty"`
    Total            int                 `json:"total"`
    Healthy          int                 `json:"healthy"`
    Unhealthy        int                 `json:"unhealthy"`
    Initial          int                 `json:"initial"`
    Draining         int                 `json:"draining"`
    Unused           int                 `json:"unused"`
    Unavailable      int                 `json:"unavailable"`
    Targets          []TargetHealthEntry `json:"targets"`
}
