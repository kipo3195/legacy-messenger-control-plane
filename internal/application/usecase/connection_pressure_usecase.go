package usecase

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/service"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
	"math"
)

type connectionPressureUsecase struct {
	ecsPort    ports.ECSPort
	ecsCfg     *configs.ECSConfig
	cloudWatch ports.CloudWatchPort
	elbPort    ports.ELBPort
	registry   *configs.ServiceRegistry
	calculator service.ConnectionPressureCalculator
}

type ConnectionPressureUsecase interface {
	GetConnectionStatus(ctx context.Context, serviceName string) (domain.ConnectionPressure, error)
}

func NewConnectionPressureUsecase(ecsPort ports.ECSPort, elbPort ports.ELBPort, cloudWatch ports.CloudWatchPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry, calculator service.ConnectionPressureCalculator) ConnectionPressureUsecase {
	return &connectionPressureUsecase{
		ecsPort:    ecsPort,
		elbPort:    elbPort,
		cloudWatch: cloudWatch,
		ecsCfg:     ecsCfg,
		registry:   registry,
		calculator: calculator,
	}
}

func (s *connectionPressureUsecase) GetConnectionStatus(ctx context.Context, serviceName string) (domain.ConnectionPressure, error) {
	return s.calculator.Calculate(ctx, serviceName)
}

func evaluateConnectionPressure(
	activeConnectionCount float64,
	connectionPerTask float64,
	runningTaskCount int,
	desiredCount int,
	minCount int,
	maxCount int,
	targetConnectionsPerTask int,
) (string, domain.ScalingRecommendation) {
	recommendation := domain.ScalingRecommendation{
		Action:                  "NONE",
		Reason:                  "connection pressure is within normal range",
		RecommendedDesiredCount: desiredCount,
	}

	if runningTaskCount <= 0 {
		recommendation.Action = "CHECK_SERVICE"
		recommendation.Reason = "running task count is zero"
		recommendation.RecommendedDesiredCount = maxInt(desiredCount, minCount)

		return "UNKNOWN", recommendation
	}

	target := float64(targetConnectionsPerTask)

	if connectionPerTask >= target*1.2 {
		recommendedDesiredCount := calculateScaleOutDesiredCount(
			activeConnectionCount,
			targetConnectionsPerTask,
			desiredCount,
			maxCount,
		)

		recommendation.Action = "SCALE_OUT"
		recommendation.Reason = "connectionPerTask is over 120% of targetConnectionsPerTask"
		recommendation.RecommendedDesiredCount = recommendedDesiredCount

		return "CRITICAL", recommendation
	}

	if connectionPerTask >= target {
		recommendedDesiredCount := calculateScaleOutDesiredCount(
			activeConnectionCount,
			targetConnectionsPerTask,
			desiredCount,
			maxCount,
		)

		recommendation.Action = "SCALE_OUT"
		recommendation.Reason = "connectionPerTask exceeds targetConnectionsPerTask"
		recommendation.RecommendedDesiredCount = recommendedDesiredCount

		return "HIGH", recommendation
	}

	if connectionPerTask >= target*0.8 {
		recommendation.Action = "WATCH"
		recommendation.Reason = "connectionPerTask is close to targetConnectionsPerTask"

		return "WARNING", recommendation
	}

	if connectionPerTask <= target*0.4 && desiredCount > minCount {
		recommendedDesiredCount := desiredCount - 1
		if recommendedDesiredCount < minCount {
			recommendedDesiredCount = minCount
		}

		recommendation.Action = "SCALE_IN_CANDIDATE"
		recommendation.Reason = "connectionPerTask is low compared to targetConnectionsPerTask"
		recommendation.RecommendedDesiredCount = recommendedDesiredCount

		return "LOW", recommendation
	}

	return "NORMAL", recommendation
}

func calculateScaleOutDesiredCount(
	activeConnectionCount float64,
	targetConnectionsPerTask int,
	currentDesiredCount int,
	maxCount int,
) int {
	if targetConnectionsPerTask <= 0 {
		return currentDesiredCount
	}

	requiredCount := int(math.Ceil(activeConnectionCount / float64(targetConnectionsPerTask)))

	// scale-out 판단이면 최소한 현재 desiredCount보다 1개는 늘리는 방향
	if requiredCount <= currentDesiredCount {
		requiredCount = currentDesiredCount + 1
	}

	if maxCount > 0 && requiredCount > maxCount {
		requiredCount = maxCount
	}

	return requiredCount
}

func roundFloat(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}
