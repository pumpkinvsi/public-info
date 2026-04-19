package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"src/backend/internal/config"
	"src/backend/internal/logging"
	"src/backend/internal/repository/postgres"
	"src/backend/internal/server"
)

const dbConnectTimeout = 10 * time.Second

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

	connectCtx, connectCancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer connectCancel()

	store, err := postgres.New(connectCtx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	slog.Info("database connection established")

	srv := server.New(cfg, store)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		slog.Error("server terminated with error", "error", err)
		os.Exit(1)
	}
}
