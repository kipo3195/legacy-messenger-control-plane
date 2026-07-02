package bootstrap

import (
	"context"
	"legacy-messenger-control-plane/internal/config"
)

type App struct {
	Config   *config.Config
	Clients  *Clients
	UseCases *UseCases
	Handlers *Handlers
	Router   *Router
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	registry, err := config.NewServiceRegistry(cfg)
	if err != nil {
		return nil, err
	}

	clients, err := NewClients(ctx, cfg)
	if err != nil {
		return nil, err
	}

	useCases := NewUseCases(clients, registry)
	handlers := NewHandlers(useCases)
	router := NewRouter(handlers)

	return &App{
		Config:   cfg,
		Clients:  clients,
		UseCases: useCases,
		Handlers: handlers,
		Router:   router,
	}, nil
}
