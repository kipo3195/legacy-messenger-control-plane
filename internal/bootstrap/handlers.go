package bootstrap

import httpadapter "legacy-messenger-control-plane/internal/adapters/http/handler"

type Handlers struct {
	ServiceObservation *httpadapter.ServiceObservationHandler
	TaskObservation    *httpadapter.TaskObservationHandler
	TaskSessionReport  *httpadapter.TaskSessionReportHandler
	TargetHealth       *httpadapter.TargetHealthHandler
	ConnectionPressure *httpadapter.ConnectionPressureHandler
	ServiceScale       *httpadapter.ServiceScaleHandler
	ServiceControl     *httpadapter.ServiceControlHandler
	ServiceEvaluation  *httpadapter.ServiceEvaluationHandler
}

func NewHandlers(useCases *UseCases) *Handlers {
	return &Handlers{
		ServiceObservation: httpadapter.NewServiceObservationHandler(
			useCases.ServiceObservationStatus,
		),
		TaskObservation: httpadapter.NewTaskObservationHandler(
			useCases.TaskObservationStatus,
		),
		TargetHealth: httpadapter.NewTargetHealthHandler(
			useCases.TargetHealth,
		),
		ConnectionPressure: httpadapter.NewConnectionPressureHandler(
			useCases.ConnectionPressure,
		),
		ServiceScale: httpadapter.NewServiceScaleHandler(
			useCases.ServiceScale,
		),
		ServiceControl: httpadapter.NewServiceControlHandler(
			useCases.ServiceControl,
		),
		ServiceEvaluation: httpadapter.NewServiceEvaluationHandler(
			useCases.ServiceEvaluation,
		),
		TaskSessionReport: httpadapter.NewTaskSessionReportHandler(
			useCases.TaskSessionReport,
		),
	}
}
