package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var enc *json.Encoder = json.NewEncoder(w)

	if v == nil {
		if _, err := w.Write([]byte("null")); err != nil {
			slog.Error("response writing failed", "error", err)
		}

		return
	}

	if err := enc.Encode(v); err != nil {
		slog.Error("response encoding failed", "error", err)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}
