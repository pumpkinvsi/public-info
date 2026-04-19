package handler

import (
	"log/slog"
	"net/http"
)

// GET /api/v1/technologies
func (h *Handler) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	techs, err := h.store.ListTechnologies(r.Context())
	if err != nil {
		slog.Error("list technologies", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, techs)
}
