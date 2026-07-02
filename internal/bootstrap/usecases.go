package bootstrap

import (
	"legacy-messenger-control-plane/internal/application"
	"legacy-messenger-control-plane/internal/config"
)

type UseCases struct {
	ServiceStatus *application.ServiceStatusUseCase
	// ServiceList        *application.ServiceListUseCase
	// ScaleService       *application.ScaleServiceUseCase
	// RedeployService    *application.RedeployServiceUseCase
	// TargetHealth       *application.TargetHealthUseCase
	// ConnectionPressure *application.ConnectionPressureUseCase
}

func NewUseCases(
	clients *Clients,
	registry *config.ServiceRegistry,
) *UseCases {
	return &UseCases{
		ServiceStatus: application.NewServiceStatusUseCase(
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
