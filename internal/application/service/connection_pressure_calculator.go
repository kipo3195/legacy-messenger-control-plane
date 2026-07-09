package service

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
	"strings"
	"time"
)

type ConnectionPressureCalculator interface {
	Calculate(ctx context.Context, serviceName string) (domain.ConnectionPressure, error)
}

type connectionPressureCalculator struct {
	ecsPort    ports.ECSPort
	ecsCfg     *configs.ECSConfig
	cloudWatch ports.CloudWatchPort
	elbPort    ports.ELBPort
	registry   *configs.ServiceRegistry
}

func NewConnectionPressureCalculator(
	ecsPort ports.ECSPort,
	elbPort ports.ELBPort,
	cloudWatch ports.CloudWatchPort,
	ecsCfg *configs.ECSConfig,
	registry *configs.ServiceRegistry,
) ConnectionPressureCalculator {
	return &connectionPressureCalculator{
		ecsPort:    ecsPort,
		ecsCfg:     ecsCfg,
		cloudWatch: cloudWatch,
		elbPort:    elbPort,
		registry:   registry,
	}
}

func (c *connectionPressureCalculator) Calculate(ctx context.Context, serviceName string) (domain.ConnectionPressure, error) {
	if serviceName == "" {
		return domain.ConnectionPressure{}, fmt.Errorf("serviceName is required")
	}

	serviceDef, err := c.registry.Find(serviceName)
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

	ecsStatus, err := c.ecsPort.DescribeService(
		ctx,
		c.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)
	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get ecs service status: %w", err)
	}

	targetGroupArn, err := c.ecsPort.GetServiceTargetGroupArn(
		ctx,
		c.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)
	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get target group arn: %w", err)
	}

	loadBalancerArn, err := c.elbPort.GetLoadBalancerArnByTargetGroupArn(
		ctx,
		targetGroupArn,
	)
	if err != nil {
		return domain.ConnectionPressure{}, fmt.Errorf("failed to get load balancer arn: %w", err)
	}

	activeConnectionCount, err := c.cloudWatch.GetALBActiveConnectionCount(
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

	pressureStatus := evaluatePressureStatus(
		connectionPerTask,
		targetConnectionsPerTask,
	)

	return domain.ConnectionPressure{
		ServiceName:    serviceName,
		ECSServiceName: serviceDef.ECSServiceName,
		ClusterName:    c.ecsCfg.ClusterName,

		ActiveConnectionCount: activeConnectionCount,
		RunningTaskCount:      runningTaskCount,
		DesiredCount:          desiredCount,

		ConnectionPerTask:        roundFloat(connectionPerTask, 2),
		TargetConnectionsPerTask: targetConnectionsPerTask,

		PressureStatus: pressureStatus,

		Metric: domain.ConnectionPressureMetric{
			Namespace:       "AWS/ApplicationELB",
			MetricName:      "ActiveConnectionCount",
			Stat:            "Sum",
			PeriodSeconds:   defaultMetricPeriodSeconds,
			LookbackMinutes: defaultMetricLookbackMinutes,
		},
	}, nil
}
