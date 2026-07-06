package bootstrap

import (
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application"
)

type UseCases struct {
	ServiceStatus *application.ServiceStatusUsecase
	// ServiceList        *application.ServiceListUseCase
	// ScaleService       *application.ScaleServiceUseCase
	// RedeployService    *application.RedeployServiceUseCase
	// TargetHealth       *application.TargetHealthUseCase
	// ConnectionPressure *application.ConnectionPressureUseCase
}

func NewUseCases(
	clients *Clients,
	registry *configs.ServiceRegistry,
) *UseCases {
	return &UseCases{
		ServiceStatus: application.NewServiceStatusUsecase(
			clients.ECS,
			registry,
		),

		// ServiceList: application.NewServiceListUseCase(
		// 	registry,
		// ),

		// ScaleService: application.NewScaleServiceUseCase(
		// 	clients.ECS,
		// 	registry,
		// ),

		// RedeployService: application.NewRedeployServiceUseCase(
		// 	clients.ECS,
		// 	registry,
		// ),

		// TargetHealth: application.NewTargetHealthUseCase(
		// 	clients.ELB,
		// 	registry,
		// ),

		// ConnectionPressure: application.NewConnectionPressureUseCase(
		// 	clients.ECS,
		// 	clients.CloudWatch,
		// 	registry,
		// ),
	}
}
