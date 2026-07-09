package application

import (
	"context"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type targetHealthUsecase struct {
	ecsPort  ports.ECSPort
	ecsCfg   *configs.ECSConfig
	elbPort  ports.ELBPort
	registry *configs.ServiceRegistry
}

type TargetHealthUsecase interface {
	GetTargetHealth(ctx context.Context, serviceName string) (*domain.TargetHealthResponse, error)
}

func NewTargetHealthUsecase(ecsPort ports.ECSPort, elbPort ports.ELBPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) TargetHealthUsecase {
	return &targetHealthUsecase{
		ecsPort:  ecsPort,
		elbPort:  elbPort,
		ecsCfg:   ecsCfg,
		registry: registry,
	}
}

func (s *targetHealthUsecase) GetTargetHealth(ctx context.Context, serviceName string) (*domain.TargetHealthResponse, error) {

	serviceDef, err := s.registry.Find(serviceName)
	if err != nil {
		return nil, err
	}

	targetGroups, err := s.ecsPort.GetServiceTargetGroups(
		ctx,
		s.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)
	if err != nil {
		return nil, err
	}

	response := &domain.TargetHealthResponse{
		ServiceName:    serviceName,
		ECSServiceName: serviceDef.ECSServiceName,
		ClusterName:    s.ecsCfg.ClusterName,
		OverallStatus:  domain.TargetHealthOverallUnknown,
		TargetGroups:   make([]domain.TargetGroupHealth, 0),
	}

	if len(targetGroups) == 0 {
		return response, nil
	}

	for _, tg := range targetGroups {
		tgHealth, err := s.elbPort.DescribeTargetHealth(
			ctx,
			tg.TargetGroupArn,
			serviceDef.LoadBalancerType,
		)
		if err != nil {
			return nil, err
		}

		response.TargetGroups = append(response.TargetGroups, *tgHealth)
		addSummary(&response.Summary, *tgHealth)
	}

	response.OverallStatus = evaluateOverallStatus(response.Summary)

	return response, nil
}

func addSummary(summary *domain.TargetHealthSummary, tg domain.TargetGroupHealth) {
	summary.Total += tg.Total
	summary.Healthy += tg.Healthy
	summary.Unhealthy += tg.Unhealthy
	summary.Initial += tg.Initial
	summary.Draining += tg.Draining
	summary.Unused += tg.Unused
	summary.Unavailable += tg.Unavailable
}

func evaluateOverallStatus(summary domain.TargetHealthSummary) domain.TargetHealthOverallStatus {
	if summary.Total == 0 {
		return domain.TargetHealthOverallUnknown
	}

	if summary.Unhealthy > 0 || summary.Unavailable > 0 {
		return domain.TargetHealthOverallDegraded
	}

	if summary.Healthy == summary.Total {
		return domain.TargetHealthOverallHealthy
	}

	if summary.Initial > 0 || summary.Draining > 0 {
		return domain.TargetHealthOverallTransitioning
	}

	return domain.TargetHealthOverallDegraded
}
