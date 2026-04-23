package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"src/backend/internal/data/bio"
	"src/backend/internal/data/contacts"
	"src/backend/internal/data/projects"
	"src/backend/internal/data/skills"
	"src/backend/internal/data/technologies"
	"src/backend/internal/email"
	"src/backend/internal/health"

	"src/backend/internal/shared/config"
	"src/backend/internal/shared/db"
	"src/backend/internal/shared/metrics"
	"src/backend/internal/shared/outbox"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 60 * time.Second
	shutdownTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, db db.Postgres, outbox outbox.Outbox) *Server {
	r := buildRouter(&db)

	return &Server{
		httpServer: &http.Server{
			Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler:      r,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		slog.Info("server listening", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	slog.Info("server stopped cleanly")
	return nil
}

func buildRouter(db *db.Postgres) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(metrics.Middleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // TODO: change to frontend
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-Request-ID"},
		ExposedHeaders: []string{"X-Request-ID"},
		MaxAge:         300,
	}))

	r.Handle("/metrics", promhttp.Handler())

	health.RegisterRoutes(r, db)
	bio.RegisterRoutes(r, db)
	contacts.RegisterRoutes(r, db)
	skills.RegisterRoutes(r, db)
	projects.RegisterRoutes(r, db)
	technologies.RegisterRoutes(r, db)
	email.RegisterRoutes(r, db)

	return r
}
