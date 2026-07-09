package bootstrap

import httpadapter "legacy-messenger-control-plane/internal/adapters/http"

type Handlers struct {
	ServiceObservation *httpadapter.ServiceObservationHandler
	TaskObservation    *httpadapter.TaskObservationHandler
	TargetHealth       *httpadapter.TargetHealthHandler
	ConnectionPressure *httpadapter.ConnectionPressureHandler
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
	}
}
