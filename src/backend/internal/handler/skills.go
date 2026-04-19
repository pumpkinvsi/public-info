package handler

import (
	"net/http"

	"src/backend/internal/model"
)

// GET /api/v1/skills
func (h *Handler) GetSkills(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	respondJSON(w, http.StatusOK, []model.Skill{})
}