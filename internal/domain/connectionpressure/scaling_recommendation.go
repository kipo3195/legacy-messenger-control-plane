package connectionpressure

type ScalingRecommendation struct {
    Action                  string `json:"action"`
    Reason                  string `json:"reason"`
    RecommendedDesiredCount int    `json:"recommendedDesiredCount"`
}
