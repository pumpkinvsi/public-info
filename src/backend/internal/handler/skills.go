package handler

import (
	"log/slog"
	"net/http"
)

// GET /api/v1/skills
func (h *Handler) GetSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := h.store.ListSkills(r.Context())
	if err != nil {
		slog.Error("list skills", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, skills)
}
