package handler

import (
	"net/http"

	"src/backend/internal/model"
)

// GET /api/v1/technologies
func (h *Handler) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	respondJSON(w, http.StatusOK, []model.Technology{})
}