package usecase

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/service"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type serviceEvaluationUsecase struct {
	ecsPort    ports.ECSPort
	ecsCfg     *configs.ECSConfig
	cloudWatch ports.CloudWatchPort
	elbPort    ports.ELBPort
	registry   *configs.ServiceRegistry
	calculator service.ConnectionPressureCalculator
}

type ServiceEvaluationUsecase interface {
	Evaluate(ctx context.Context, serviceName string) (domain.ScalingEvaluation, error)
}

func NewServiceEvaluationUsecase(ecsPort ports.ECSPort, elbPort ports.ELBPort, cloudWatch ports.CloudWatchPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry, calculator service.ConnectionPressureCalculator) ServiceEvaluationUsecase {
	return &serviceEvaluationUsecase{
		ecsPort:    ecsPort,
		elbPort:    elbPort,
		cloudWatch: cloudWatch,
		ecsCfg:     ecsCfg,
		registry:   registry,
		calculator: calculator,
	}
}
func (u *serviceEvaluationUsecase) Evaluate(ctx context.Context, serviceName string) (domain.ScalingEvaluation, error) {
	if serviceName == "" {
		return domain.ScalingEvaluation{}, fmt.Errorf("serviceName is required")
	}

	serviceDef, err := u.registry.Find(serviceName)
	if err != nil {
		return domain.ScalingEvaluation{}, fmt.Errorf("service not found: %s", serviceName)
	}

	pressure, err := u.calculator.Calculate(ctx, serviceName)
	if err != nil {
		return domain.ScalingEvaluation{}, err
	}

	target := serviceDef.TargetConnectionsPerTask
	if target <= 0 {
		target = 1500
	}

	scaleOutThreshold := int(float64(target) * 0.8)
	scaleInThreshold := int(float64(target) * 0.3)

	recommendedDesired := calculateRecommendedDesiredCount(
		pressure.ActiveConnectionCount,
		target,
		serviceDef.MinCount,
		serviceDef.MaxCount,
	)

	action := domain.ScalingActionKeep
	reason := "connection pressure is within threshold"

	if !serviceDef.Scalable {
		action = domain.ScalingActionNotScalable
		reason = "service is not scalable"
		recommendedDesired = pressure.DesiredCount

	} else if pressure.ConnectionPerTask >= float64(scaleOutThreshold) {
		if pressure.DesiredCount >= serviceDef.MaxCount {
			action = domain.ScalingActionKeep
			reason = "connectionPerTask is above scale-out threshold, but desiredCount already reached maxCount"
			recommendedDesired = pressure.DesiredCount
		} else {
			action = domain.ScalingActionScaleOut
			reason = "connectionPerTask is above scale-out threshold"

			if recommendedDesired <= pressure.DesiredCount {
				recommendedDesired = pressure.DesiredCount + 1
			}

			if recommendedDesired > serviceDef.MaxCount {
				recommendedDesired = serviceDef.MaxCount
			}
		}

	} else if pressure.ConnectionPerTask <= float64(scaleInThreshold) {
		if pressure.DesiredCount <= serviceDef.MinCount {
			action = domain.ScalingActionKeep
			reason = "connectionPerTask is below scale-in threshold, but desiredCount already reached minCount"
			recommendedDesired = pressure.DesiredCount
		} else {
			action = domain.ScalingActionScaleIn
			reason = "connectionPerTask is below scale-in threshold"

			if recommendedDesired >= pressure.DesiredCount {
				recommendedDesired = pressure.DesiredCount - 1
			}

			if recommendedDesired < serviceDef.MinCount {
				recommendedDesired = serviceDef.MinCount
			}
		}

	} else {
		recommendedDesired = pressure.DesiredCount
	}

	return domain.ScalingEvaluation{
		ServiceName:    serviceName,
		ECSServiceName: serviceDef.ECSServiceName,
		ClusterName:    u.ecsCfg.ClusterName,
		Action:         action,
		Reason:         reason,
		Current: domain.ScalingCurrentStatus{
			ActiveConnectionCount: pressure.ActiveConnectionCount,
			RunningTaskCount:      pressure.RunningTaskCount,
			DesiredCount:          pressure.DesiredCount,
			ConnectionPerTask:     pressure.ConnectionPerTask,
		},
		Policy: domain.ScalingPolicyStatus{
			TargetConnectionsPerTask: target,
			ScaleOutThreshold:        scaleOutThreshold,
			ScaleInThreshold:         scaleInThreshold,
			MinCount:                 serviceDef.MinCount,
			MaxCount:                 serviceDef.MaxCount,
		},
		Recommendation: domain.ScalingRecommendationStatus{
			CurrentDesiredCount:     pressure.DesiredCount,
			RecommendedDesiredCount: recommendedDesired,
		},
	}, nil
}
