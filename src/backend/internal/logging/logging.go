package logging

import (
	"context"
	"log/slog"
	"os"
)

const lokiPushPath = "/loki/api/v1/push"

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
			_ = h.Handle(ctx, r)
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

func Setup(level slog.Level, lokiAddr string, labels map[string]string) {
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	lokiHandler := newLokiHandler(lokiAddr+lokiPushPath, labels)

	multi := newMultiHandler(stdoutHandler, lokiHandler)

	slog.SetDefault(slog.New(multi))
}
