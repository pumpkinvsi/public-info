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

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

type LokiHandler struct {
	endpoint string
	client   *http.Client
	labels   map[string]string
	attrs    []slog.Attr
	groups   []string
}

func newLokiHandler(endpoint string, labels map[string]string) *LokiHandler {
	return &LokiHandler{
		endpoint: endpoint,
		labels:   labels,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (h *LokiHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

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
		return nil
	}

	resp, err := h.client.Post(h.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	return nil
}

func (h *LokiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := h.clone()
	clone.attrs = append(clone.attrs, attrs...)
	return clone
}

func (h *LokiHandler) WithGroup(name string) slog.Handler {
	clone := h.clone()
	clone.groups = append(clone.groups, name)
	return clone
}

func (h *LokiHandler) format(r slog.Record) string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "level=%s msg=%q", r.Level, r.Message)

	for _, a := range h.attrs {
		fmt.Fprintf(buf, " %s=%v", h.qualifyKey(a.Key), a.Value)
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(buf, " %s=%v", h.qualifyKey(a.Key), a.Value)
		return true
	})

	return buf.String()
}

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
		client:   h.client,
		labels:   labels,
		attrs:    append([]slog.Attr{}, h.attrs...),
		groups:   append([]string{}, h.groups...),
	}
}
