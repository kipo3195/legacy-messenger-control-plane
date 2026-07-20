package bootstrap

import (
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/application/service"
	"legacy-messenger-control-plane/internal/application/usecase"
)

type UseCases struct {
	ServiceObservationStatus usecase.ServiceObservationUsecase
	TaskObservationStatus    usecase.TaskObservationUsecase
	TaskSessionReport        usecase.TaskSessionReportUsecase
	ServiceScale             usecase.ServiceScaleUsecase
	ServiceControl           usecase.ServiceControlUsecase
	TargetHealth             usecase.TargetHealthUsecase
	ConnectionPressure       usecase.ConnectionPressureUsecase
	ServiceEvaluation        usecase.ServiceEvaluationUsecase
	AutoScale                usecase.SessionAutoScalingUsecase
	ScaleInUsecase           usecase.ScaleInUsecase
}

func NewUseCases(clients *Clients, cfg *configs.Config, registry *configs.ServiceRegistry) *UseCases {

	connectionPressureCalculator := service.NewConnectionPressureCalculator(
		clients.ECS,
		clients.ELB,
		clients.CloudWatch,
		cfg.ECS,
		registry,
	)

	scalingPolicy := usecase.NewScalingPolicy()
	scaleInCoordinator := usecase.NewScaleInCoordinator()

	return &UseCases{
		ServiceObservationStatus: usecase.NewServiceObservationUsecase(
			clients.ECS,
			cfg.ECS,
			registry,
		),
		TaskObservationStatus: usecase.NewTaskObservationUsecase(
			clients.ECS,
			cfg.ECS,
			registry,
		),
		ServiceScale: usecase.NewServiceScaleUsecase(
			clients.ECS,
			cfg.ECS,
			registry,
		),

		TargetHealth: usecase.NewTargetHealthUsecase(
			clients.ECS,
			clients.ELB,
			cfg.ECS,
			registry,
		),

		ConnectionPressure: usecase.NewConnectionPressureUsecase(
			clients.ECS,
			clients.ELB,
			clients.CloudWatch,
			cfg.ECS,
			registry,
			connectionPressureCalculator,
		),

		ServiceControl: usecase.NewServiceControlUsecase(
			clients.ECS,
			clients.ELB,
			clients.CloudWatch,
			cfg.ECS,
			registry,
		),

		ServiceEvaluation: usecase.NewServiceEvaluationUsecase(
			clients.ECS,
			clients.ELB,
			clients.CloudWatch,
			cfg.ECS,
			registry,
			connectionPressureCalculator,
		),

		TaskSessionReport: usecase.NewTaskSessionReportUsecase(
			clients.TaskSession,
		),

		AutoScale: usecase.NewSessionAutoScalingUsecase(
			clients.TaskSession,
			clients.ECS,
			registry,
			cfg.ECS,
			cfg.AutoScale,
			scalingPolicy,
		),

		ScaleInUsecase: usecase.NewScaleInUsecase(
			clients.TaskSession,
			clients.ECS,
			scalingPolicy,
			scaleInCoordinator,
		),
	}
}
