package bootstrap

import (
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/usecase"
)

type UseCases struct {
	ServiceObservationStatus usecase.ServiceObservationUsecase
	TaskObservationStatus    usecase.TaskObservationUsecase
	ServiceScale             usecase.ServiceScaleUsecase
	// RedeployService    *application.RedeployServiceUseCase
	TargetHealth       usecase.TargetHealthUsecase
	ConnectionPressure usecase.ConnectionPressureUsecase
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

		// RedeployService: application.NewRedeployServiceUseCase(
		// 	clients.ECS,
		// 	registry,
		// ),

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
	}
}
