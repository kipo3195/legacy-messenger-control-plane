package bootstrap

import httpadapter "legacy-messenger-control-plane/internal/adapters/http"

type Handlers struct {
	Service *httpadapter.ServiceHandler
}

func NewHandlers(useCases *UseCases) *Handlers {
	return &Handlers{
		Service: httpadapter.NewServiceHandler(
			useCases.ServiceStatus,
		),
	}
}
