package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"src/backend/internal/config"
	"src/backend/internal/repository"
)

type Handler struct {
	cfg   *config.Config
	store repository.Store
}

func New(cfg *config.Config, store repository.Store) *Handler {
	return &Handler{cfg: cfg, store: store}
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("response encoding failed", "error", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
