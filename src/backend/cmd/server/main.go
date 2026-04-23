package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"src/backend/internal/server"
	"src/backend/internal/shared/config"
	"src/backend/internal/shared/db"
	"src/backend/internal/shared/kafka"
	"src/backend/internal/shared/logging"
	"src/backend/internal/shared/outbox"
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

	store, err := db.New(connectCtx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	slog.Info("database connection established")

	producer, err := kafka.New(cfg.Kafka)
	if err != nil {
		slog.Error("failed to create kafka producer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			slog.Error("kafka producer close", "error", err)
		}
	}()

	outboxRepository := outbox.NewRepository(store)

	pool := outbox.NewWorkerPool(
		cfg.Kafka.WorkerCount,
		cfg.Kafka.PollInterval,
		outboxRepository,
		producer,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancel()

	go pool.Run(ctx)

	srv := server.New(cfg, *store, outboxRepository)

	if err := srv.Run(ctx); err != nil {
		slog.Error("server terminated with error", "error", err)
		os.Exit(1)
	}
}
