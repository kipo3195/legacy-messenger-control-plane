package bootstrap

import httpadapter "legacy-messenger-control-plane/internal/adapters/http"

type Handlers struct {
	Service *httpadapter.ServiceObservationHandler
}

func NewHandlers(useCases *UseCases) *Handlers {
	return &Handlers{
		Service: httpadapter.NewServiceObservationHandler(
			useCases.ServiceObservationStatus,
		),
	}
}
