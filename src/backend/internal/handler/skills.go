package handler

import (
	"log/slog"
	"net/http"
)

// GET /api/v1/skills
func (h *Handler) GetSkills(w http.ResponseWriter, r *http.Request) {
	groups, err := h.store.ListProjectsGrouped(r.Context())
	if err != nil {
		slog.Error("list projects grouped", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, groups)
}
