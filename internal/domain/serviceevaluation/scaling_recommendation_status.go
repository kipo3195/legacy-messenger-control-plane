package serviceevaluation

type ScalingRecommendationStatus struct {
	CurrentDesiredCount     int `json:"currentDesiredCount"`
	RecommendedDesiredCount int `json:"recommendedDesiredCount"`
}
