package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"src/email-sender/internal/config"
	"src/email-sender/internal/email"
	"src/email-sender/internal/kafka"
	"src/email-sender/internal/logging"
	"src/email-sender/internal/metrics"
	"src/email-sender/internal/smtp"
)

const dbConnectTimeout = 10 * time.Second

func main() {

	logging.Setup(
		slog.LevelInfo,
		"http://loki:3100",
		map[string]string{
			"app": "email-sender",
		},
	)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}
	m := metrics.New()
	metricsAddr := fmt.Sprintf(":%d", cfg.Metrics.Port)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		slog.Info("metrics server listening", "addr", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, mux); err != nil && err != http.ErrServerClosed {
			slog.Error("metrics server error", "error", err)
		}
	}()

	smtpClient := smtp.NewSmtpClient(&cfg.Smtp)

	jobs := make(chan *email.Email, 100)

	pool := email.NewPool(jobs, smtpClient, m)
	pool.Start()

	consumer, err := kafka.New(cfg.Kafka, m)
	if err != nil {
		slog.Error("create kafka consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		consumer.Listen(ctx, jobs)
		close(jobs)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	slog.Info("shutdown signal received", "signal", sig.String())

	cancel()
	pool.Wait()

	slog.Info("email-sender stopped cleanly")
}
