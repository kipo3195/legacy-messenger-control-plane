package main

import (
	"context"
	"legacy-messenger-control-plane/internal/bootstrap"
	"log"
)

func main() {
	ctx := context.Background()

	app, err := bootstrap.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	addr := ":" + app.Config.Server.Port

	if err := app.Router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
