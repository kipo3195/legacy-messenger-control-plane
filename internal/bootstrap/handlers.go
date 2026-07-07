package bootstrap

import httpadapter "legacy-messenger-control-plane/internal/adapters/http"

type Handlers struct {
	ServiceObservation *httpadapter.ServiceObservationHandler
	TaskObservation    *httpadapter.TaskObservationHandler
}

func NewHandlers(useCases *UseCases) *Handlers {
	return &Handlers{
		ServiceObservation: httpadapter.NewServiceObservationHandler(
			useCases.ServiceObservationStatus,
		),
		TaskObservation: httpadapter.NewTaskObservationHandler(
			useCases.TaskObservationStatus,
		),
	}
}
