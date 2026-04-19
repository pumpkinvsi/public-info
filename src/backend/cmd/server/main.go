package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"src/backend/internal/config"
	"src/backend/internal/logging"
	"src/backend/internal/server"
)

func main() {

	logging.Setup(
		slog.LevelInfo,
		"http://loki:3100",
		map[string]string{
			"app": "backend",
		},
	)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	srv := server.New(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		slog.Error("server terminated with error", "error", err)
		os.Exit(1)
	}
}
