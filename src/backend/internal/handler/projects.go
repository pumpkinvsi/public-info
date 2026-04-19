package handler

import (
	"net/http"

	"src/backend/internal/model"
)

// GET /api/v1/projects
func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	respondJSON(w, http.StatusOK, []model.ProjectGroup{})
}