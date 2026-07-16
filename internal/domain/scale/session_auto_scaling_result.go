package scale

import (
	"legacy-messenger-control-plane/internal/domain/serviceevaluation"
)

type SessionAutoScalingResult struct {
	ServiceName string `json:"serviceName"`

	TotalSessionCount       int `json:"totalSessionCount"`
	CurrentDesiredCount     int `json:"currentDesiredCount"`
	RecommendedDesiredCount int `json:"recommendedDesiredCount"`

	Action   serviceevaluation.ScalingAction `json:"action"`
	Executed bool                            `json:"executed"`
	Reason   string                          `json:"reason"`

	ECSState interface{} `json:"ecsState,omitempty"`
}
