package bootstrap

import (
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/usecase"
)

type UseCases struct {
	ServiceObservationStatus usecase.ServiceObservationUsecase
	TaskObservationStatus    usecase.TaskObservationUsecase
	ServiceScale             usecase.ServiceScaleUsecase
	ServiceControl           usecase.ServiceControlUsecase
	TargetHealth             usecase.TargetHealthUsecase
	ConnectionPressure       usecase.ConnectionPressureUsecase
}

func NewUseCases(clients *Clients, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) *UseCases {
	return &UseCases{
		ServiceObservationStatus: usecase.NewServiceObservationUsecase(
			clients.ECS,
			ecsCfg,
			registry,
		),
		TaskObservationStatus: usecase.NewTaskObservationUsecase(
			clients.ECS,
			ecsCfg,
			registry,
		),
		ServiceScale: usecase.NewServiceScaleUsecase(
			clients.ECS,
			ecsCfg,
			registry,
		),

		TargetHealth: usecase.NewTargetHealthUsecase(
			clients.ECS,
			clients.ELB,
			ecsCfg,
			registry,
		),

		ConnectionPressure: usecase.NewConnectionPressureUsecase(
			clients.ECS,
			clients.ELB,
			clients.CloudWatch,
			ecsCfg,
			registry,
		),

		ServiceControl: usecase.NewServiceControlUsecase(
			clients.ECS,
			clients.ELB,
			clients.CloudWatch,
			ecsCfg,
			registry,
		),
	}
}
