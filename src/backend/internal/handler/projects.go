package handler

import (
	"net/http"
)

// GET /api/v1/projects
func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.ListProjectsGrouped(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, projects)
}
