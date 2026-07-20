package serviceevaluation

type ScalingAction string

const (
	ScalingActionScaleOut         ScalingAction = "SCALE_OUT"
	ScalingActionScaleIn          ScalingAction = "SCALE_IN"
	ScalingActionScaleInCandidate ScalingAction = "SCALE_IN_CANDIDATE"
	ScalingActionKeep             ScalingAction = "KEEP"
	ScalingActionNotScalable      ScalingAction = "NOT_SCALABLE"
	ScaleActionSkip               ScalingAction = "SCALE_ACTION_SKIP"
)

type ScalingEvaluation struct {
	ServiceName    string        `json:"serviceName"`
	ECSServiceName string        `json:"ecsServiceName"`
	ClusterName    string        `json:"clusterName"`
	Action         ScalingAction `json:"action"`
	Reason         string        `json:"reason"`

	Current ScalingCurrentStatus `json:"current"`
	Policy  ScalingPolicyStatus  `json:"policy"`

	Recommendation ScalingRecommendationStatus `json:"recommendation"`
}
