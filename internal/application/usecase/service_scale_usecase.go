package usecase

import (
	"context"
	"errors"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/command"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type serviceScaleUsecase struct {
	ecsPort  ports.ECSPort
	ecsCfg   *configs.ECSConfig
	registry *configs.ServiceRegistry
}

var (
	ErrServiceNotFound        = errors.New("service not found")
	ErrServiceNotScalable     = errors.New("service is not scalable")
	ErrDesiredCountOutOfRange = errors.New("desiredCount is out of range")
	ErrInvalidDesiredCount    = errors.New("desiredCount must be greater than or equal to 0")
)

type ServiceScaleUsecase interface {
	UpdateServiceDesirdCount(ctx context.Context, cmd command.ScaleServiceCommand) (domain.ServiceScaleResult, error)
}

func NewServiceScaleUsecase(ecsPort ports.ECSPort, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) ServiceScaleUsecase {
	return &serviceScaleUsecase{
		ecsPort:  ecsPort,
		ecsCfg:   ecsCfg,
		registry: registry,
	}
}

func (u *serviceScaleUsecase) UpdateServiceDesirdCount(ctx context.Context, cmd command.ScaleServiceCommand) (domain.ServiceScaleResult, error) {

	if cmd.DesiredCount < 0 {
		return domain.ServiceScaleResult{}, ErrInvalidDesiredCount
	}

	serviceDef, err := u.registry.Find(cmd.ServiceName)
	if err != nil {
		return domain.ServiceScaleResult{}, fmt.Errorf("%w: %s", ErrServiceNotFound, cmd.ServiceName)
	}

	if !serviceDef.Scalable {
		return domain.ServiceScaleResult{}, fmt.Errorf("%w: %s", ErrServiceNotScalable, cmd.ServiceName)
	}

	if cmd.DesiredCount < serviceDef.MinCount || cmd.DesiredCount > serviceDef.MaxCount {
		return domain.ServiceScaleResult{}, fmt.Errorf(
			"%w: desiredCount=%d, min=%d, max=%d",
			ErrDesiredCountOutOfRange,
			cmd.DesiredCount,
			serviceDef.MinCount,
			serviceDef.MaxCount,
		)
	}

	before, err := u.ecsPort.GetServiceControlState(
		ctx,
		u.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)
	if err != nil {
		return domain.ServiceScaleResult{}, err
	}

	if int(before.DesiredCount) == cmd.DesiredCount {
		return domain.ServiceScaleResult{
			ServiceName:          cmd.ServiceName,
			ECSServiceName:       serviceDef.ECSServiceName,
			ClusterName:          u.ecsCfg.ClusterName,
			PreviousDesiredCount: before.DesiredCount,
			DesiredCount:         before.DesiredCount,
			RunningCount:         before.RunningCount,
			PendingCount:         before.PendingCount,
			Status:               "NOOP",
			Message:              "desiredCount is already set",
		}, nil
	}

	after, err := u.ecsPort.UpdateServiceDesiredCount(
		ctx,
		u.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
		cmd.DesiredCount,
	)
	if err != nil {
		return domain.ServiceScaleResult{}, err
	}

	return domain.ServiceScaleResult{
		ServiceName:          cmd.ServiceName,
		ECSServiceName:       serviceDef.ECSServiceName,
		ClusterName:          u.ecsCfg.ClusterName,
		PreviousDesiredCount: before.DesiredCount,
		DesiredCount:         after.DesiredCount,
		RunningCount:         after.RunningCount,
		PendingCount:         after.PendingCount,
		Status:               "SCALING_REQUESTED",
		Message:              "scale request accepted",
	}, nil
}
