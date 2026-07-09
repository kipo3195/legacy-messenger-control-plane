package usecase

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/command"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type serviceControlUsecase struct {
	ecsPort  ports.ECSPort
	ecsCfg   *configs.ECSConfig
	elbPort  ports.ELBPort
	registry *configs.ServiceRegistry
}

type ServiceControlUsecase interface {
	ServiceRedeploy(ctx context.Context, cmd command.ServiceRedeployCommand) (domain.ServiceRedeployResult, error)
}

func NewServiceControlUsecase(ecsPort ports.ECSPort, elbPort ports.ELBPort, cloudWatch ports.CloudWatchPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) ServiceControlUsecase {
	return &serviceControlUsecase{
		ecsPort:  ecsPort,
		elbPort:  elbPort,
		ecsCfg:   ecsCfg,
		registry: registry,
	}
}

func (u *serviceControlUsecase) ServiceRedeploy(ctx context.Context, cmd command.ServiceRedeployCommand) (domain.ServiceRedeployResult, error) {

	serviceDef, err := u.registry.Find(cmd.ServiceName)
	if err != nil {
		return domain.ServiceRedeployResult{}, fmt.Errorf("service not registered: %s", cmd.ServiceName)
	}

	if serviceDef.ECSServiceName == "" {
		return domain.ServiceRedeployResult{}, fmt.Errorf("ecs service name is empty: %s", cmd.ServiceName)
	}

	before, err := u.ecsPort.DescribeService(ctx, u.ecsCfg.ClusterName, serviceDef.ECSServiceName)
	if err != nil {
		return domain.ServiceRedeployResult{}, err
	}

	if before.DesiredCount == 0 {
		return domain.ServiceRedeployResult{}, fmt.Errorf("cannot redeploy service with desiredCount 0: %s", cmd.ServiceName)
	}

	if hasDeploymentInProgress(before.Deployments) {
		return domain.ServiceRedeployResult{}, fmt.Errorf("service deployment is already in progress: %s", cmd.ServiceName)
	}

	result, err := u.ecsPort.ForceNewDeployment(ctx, u.ecsCfg.ClusterName, serviceDef.ECSServiceName)
	if err != nil {
		return domain.ServiceRedeployResult{}, err
	}

	result.ServiceName = cmd.ServiceName
	result.ECSServiceName = serviceDef.ECSServiceName
	result.ClusterName = u.ecsCfg.ClusterName
	result.Action = "REDEPLOY"
	result.Status = "STARTED"
	result.Reason = cmd.Reason
	result.Message = "force new deployment requested"

	fmt.Printf("ServiceRedeploy end.. service : %s", cmd.ServiceName)
	return result, nil
}

func hasDeploymentInProgress(deployments []domain.DeploymentStatus) bool {
	for _, d := range deployments {
		if d.Status == "PRIMARY" && d.RolloutState == "IN_PROGRESS" {
			return true
		}

		if d.Status == "ACTIVE" {
			return true
		}
	}

	return false
}
