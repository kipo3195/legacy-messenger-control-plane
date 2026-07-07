package bootstrap

import (
	"context"
	"legacy-messenger-control-plane/configs"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config   *configs.Config
	Clients  *Clients
	UseCases *UseCases
	Handlers *Handlers
	Router   *gin.Engine
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := configs.Load()
	if err != nil {
		return nil, err
	}

	serviceRegistry, err := configs.NewServiceRegistry(cfg)
	if err != nil {
		return nil, err
	}

	clients, err := NewClients(ctx, cfg)
	if err != nil {
		return nil, err
	}

	useCases := NewUseCases(clients, cfg.ECS, serviceRegistry)
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
