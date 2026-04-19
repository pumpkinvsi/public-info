package handler

import (
	"net/http"

	"src/backend/internal/model"
)

// GET /api/v1/contacts
func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	respondJSON(w, http.StatusOK, model.Contacts{Contacts: []model.Contact{}})
}