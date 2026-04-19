package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// lokiPushRequest is the JSON body accepted by POST /loki/api/v1/push.
// Reference: https://grafana.com/docs/loki/latest/reference/loki-http-api/#push-log-entries-to-loki
type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

// lokiStream groups log lines that share the same label set.
type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"` // [nanosecond unix timestamp string, log line]
}

// LokiHandler is a slog.Handler that ships every log record to Loki via its
// HTTP push API. It composes with any other slog.Handler through MultiHandler.
//
// Each log record is sent as a separate HTTP request. This is intentional for
// a low-traffic personal project. A production service should use batching.
type LokiHandler struct {
	endpoint string // full URL: http://loki:3100/loki/api/v1/push
	client   *http.Client
	labels   map[string]string // static Loki stream labels
	attrs    []slog.Attr       // accumulated via WithAttrs
	groups   []string          // accumulated via WithGroup
}

// newLokiHandler constructs a LokiHandler.
// labels become the Loki stream selector (e.g. {app="backend", env="prod"}).
func newLokiHandler(endpoint string, labels map[string]string) *LokiHandler {
	return &LokiHandler{
		endpoint: endpoint,
		labels:   labels,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// Enabled reports whether the handler handles records at the given level.
// LokiHandler ships all levels — filtering is the caller's responsibility.
func (h *LokiHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle ships the record to Loki. Errors are silently dropped to avoid
// log-induced panics disrupting the application.
func (h *LokiHandler) Handle(_ context.Context, r slog.Record) error {
	line := h.format(r)
	ts := strconv.FormatInt(r.Time.UnixNano(), 10)

	payload := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: h.labels,
				Values: [][2]string{{ts, line}},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil // never disrupt caller
	}

	resp, err := h.client.Post(h.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil // Loki unavailability must not crash the app
	}
	defer resp.Body.Close()

	return nil
}

// WithAttrs returns a new handler with the given attributes appended.
func (h *LokiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := h.clone()
	clone.attrs = append(clone.attrs, attrs...)
	return clone
}

// WithGroup returns a new handler with the given group name pushed onto the stack.
func (h *LokiHandler) WithGroup(name string) slog.Handler {
	clone := h.clone()
	clone.groups = append(clone.groups, name)
	return clone
}

// format renders a slog.Record into a plain-text log line that Loki will index.
// Attributes accumulated via WithAttrs and those in the record itself are
// appended as key=value pairs.
func (h *LokiHandler) format(r slog.Record) string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "level=%s msg=%q", r.Level, r.Message)

	// Attrs accumulated through WithAttrs
	for _, a := range h.attrs {
		fmt.Fprintf(buf, " %s=%v", h.qualifyKey(a.Key), a.Value)
	}

	// Attrs attached to this specific record
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(buf, " %s=%v", h.qualifyKey(a.Key), a.Value)
		return true
	})

	return buf.String()
}

// qualifyKey prepends group names to an attribute key, matching slog conventions.
func (h *LokiHandler) qualifyKey(key string) string {
	for i := len(h.groups) - 1; i >= 0; i-- {
		key = h.groups[i] + "." + key
	}
	return key
}

func (h *LokiHandler) clone() *LokiHandler {
	labels := make(map[string]string, len(h.labels))
	for k, v := range h.labels {
		labels[k] = v
	}
	return &LokiHandler{
		endpoint: h.endpoint,
		client:   h.client, // safe to share — http.Client is concurrent-safe
		labels:   labels,
		attrs:    append([]slog.Attr{}, h.attrs...),
		groups:   append([]string{}, h.groups...),
	}
}
