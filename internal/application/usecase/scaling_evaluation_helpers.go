package usecase

import "math"

func calculateRecommendedDesiredCount(
	activeConnections float64,
	targetPerTask int,
	minCount int,
	maxCount int,
) int {
	if targetPerTask <= 0 {
		targetPerTask = 1500
	}

	recommended := int(math.Ceil(activeConnections / float64(targetPerTask)))

	if recommended < minCount {
		recommended = minCount
	}

	if recommended > maxCount {
		recommended = maxCount
	}

	return recommended
}
