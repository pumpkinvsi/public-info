package handler

import (
	"net/http"

	"src/backend/internal/model"
)

// GET /api/v1/bio
func (h *Handler) GetBio(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	respondJSON(w, http.StatusOK, model.Bio{})
}