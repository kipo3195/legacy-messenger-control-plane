package application

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
	"math"
	"strings"
	"time"
)

type connectionPressureUsecase struct {
	ecsPort    ports.ECSPort
	ecsCfg     *configs.ECSConfig
	cloudWatch ports.CloudWatchPort
	elbPort    ports.ELBPort
	registry   *configs.ServiceRegistry
}

type ConnectionPressureUsecase interface {
	GetConnectionStatus(ctx context.Context, serviceName string) (domain.ConnectionPressure, error)
}

func NewConnectionPressureUsecase(ecsPort ports.ECSPort, elbPort ports.ELBPort, cloudWatch ports.CloudWatchPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) ConnectionPressureUsecase {
	return &connectionPressureUsecase{
		ecsPort:    ecsPort,
		elbPort:    elbPort,
		cloudWatch: cloudWatch,
		ecsCfg:     ecsCfg,
		registry:   registry,
	}
}

const (
	defaultTargetConnectionsPerTask = 1500
	defaultMetricPeriodSeconds      = int32(60)
	defaultMetricLookbackMinutes    = 5
)

func (s *connectionPressureUsecase) GetConnectionStatus(ctx context.Context, serviceName string) (domain.ConnectionPressure, error) {
	if serviceName == "" {
		return domain.ConnectionPressure{}, fmt.Errorf("serviceName is required")
	}

	serviceDef, err := s.registry.Find(serviceName)
	if err != nil {
		return domain.ConnectionPressure{}, err
	}

	if serviceDef.ECSServiceName == "" {
		return domain.ConnectionPressure{}, fmt.Errorf("ecsServiceName is empty for service: %s", serviceName)
	}

	if strings.ToLower(serviceDef.LoadBalancerType) != "alb" {
		return domain.ConnectionPressure{}, fmt.Errorf(
			"connection-pressure currently supports alb only. serviceName=%s, loadBalancerType=%s",
			serviceName,
			serviceDef.LoadBalancerType,
		)
	}

	targetConnectionsPerTask := serviceDef.TargetConnectionsPerTask
	if targetConnectionsPerTask <= 0 {
		targetConnectionsPerTask = defaultTargetConnectionsPerTask
	}

	ecsStatus, err := s.ecsPort.DescribeService(
		ctx,
		s.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)

	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get ecs service status: %w", err)
	}

	targetGroupArn, err := s.ecsPort.GetServiceTargetGroupArn(
		ctx,
		s.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)
	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get target group arn: %w", err)
	}

	loadBalancerArn, err := s.elbPort.GetLoadBalancerArnByTargetGroupArn(
		ctx,
		targetGroupArn,
	)

	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get load balancer arn: %w", err)
	}

	activeConnectionCount, err := s.cloudWatch.GetALBActiveConnectionCount(
		ctx,
		loadBalancerArn,
		int32(defaultMetricPeriodSeconds),
		time.Duration(defaultMetricLookbackMinutes)*time.Minute,
	)

	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get active connection count: %w", err)
	}

	runningTaskCount := int(ecsStatus.RunningCount)
	desiredCount := int(ecsStatus.DesiredCount)

	connectionPerTask := 0.0
	if runningTaskCount > 0 {
		connectionPerTask = activeConnectionCount / float64(runningTaskCount)
	}

	pressureStatus, recommendation := evaluateConnectionPressure(
		activeConnectionCount,
		connectionPerTask,
		runningTaskCount,
		desiredCount,
		serviceDef.MinCount,
		serviceDef.MaxCount,
		targetConnectionsPerTask,
	)

	return domain.ConnectionPressure{
		ServiceName:    serviceName,
		ECSServiceName: serviceDef.ECSServiceName,
		ClusterName:    s.ecsCfg.ClusterName,

		ActiveConnectionCount: activeConnectionCount,
		RunningTaskCount:      runningTaskCount,
		DesiredCount:          desiredCount,

		ConnectionPerTask:        roundFloat(connectionPerTask, 2),
		TargetConnectionsPerTask: targetConnectionsPerTask,

		PressureStatus:        pressureStatus,
		ScalingRecommendation: recommendation,

		Metric: domain.ConnectionPressureMetric{
			Namespace:       "AWS/ApplicationELB",
			MetricName:      "ActiveConnectionCount",
			Stat:            "Average",
			PeriodSeconds:   defaultMetricPeriodSeconds,
			LookbackMinutes: defaultMetricLookbackMinutes,
		},
	}, nil
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
