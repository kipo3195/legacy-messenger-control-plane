package bootstrap

import (
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application"
)

type UseCases struct {
	ServiceObservationStatus application.ServiceObservationUsecase
	TaskObservationStatus    application.TaskObservationUsecase
	// ScaleService       *application.ScaleServiceUseCase
	// RedeployService    *application.RedeployServiceUseCase
	TargetHealth       application.TargetHealthUsecase
	ConnectionPressure application.ConnectionPressureUsecase
}

func NewUseCases(clients *Clients, ecsCfg *configs.ECSConfig, registry *configs.ServiceRegistry) *UseCases {
	return &UseCases{
		ServiceObservationStatus: application.NewServiceObservationUsecase(
			clients.ECS,
			ecsCfg,
			registry,
		),
		TaskObservationStatus: application.NewTaskObservationUsecase(
			clients.ECS,
			ecsCfg,
			registry,
		),
		// ScaleService: application.NewScaleServiceUseCase(
		// 	clients.ECS,
		// 	registry,
		// ),

		// RedeployService: application.NewRedeployServiceUseCase(
		// 	clients.ECS,
		// 	registry,
		// ),

		TargetHealth: application.NewTargetHealthUsecase(
			clients.ECS,
			clients.ELB,
			ecsCfg,
			registry,
		),

		ConnectionPressure: application.NewConnectionPressureUsecase(
			clients.ECS,
			clients.ELB,
			clients.CloudWatch,
			ecsCfg,
			registry,
		),
	}
}
