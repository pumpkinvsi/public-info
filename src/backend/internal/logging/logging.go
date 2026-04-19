package logging

import (
	"context"
	"log/slog"
	"os"
)

const lokiPushPath = "/loki/api/v1/push"

// MultiHandler fans a single slog record out to multiple slog.Handler
// implementations. All handlers receive every record regardless of individual
// Enabled results — the top-level logger's level gate is the single filter.
type MultiHandler struct {
	handlers []slog.Handler
}

func newMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			_ = h.Handle(ctx, r) // errors are best-effort
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return newMultiHandler(handlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return newMultiHandler(handlers...)
}

// Setup initialises the global slog logger with two destinations:
//   - stdout: JSON format, consumed by the Docker log driver / Loki agent
//   - Loki:   direct HTTP push for structured stream queries in Grafana
//
// lokiAddr is the base URL of the Loki instance (e.g. "http://loki:3100").
// labels are the static Loki stream labels attached to every log line
// (e.g. map[string]string{"app": "backend"}).
func Setup(level slog.Level, lokiAddr string, labels map[string]string) {
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	lokiHandler := newLokiHandler(lokiAddr+lokiPushPath, labels)

	multi := newMultiHandler(stdoutHandler, lokiHandler)

	slog.SetDefault(slog.New(multi))
}
