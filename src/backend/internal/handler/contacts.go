package handler

import (
	"log/slog"
	"net/http"

	"src/backend/internal/model"
)

// GET /api/v1/contacts
func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.store.ListContacts(r.Context())
	if err != nil {
		slog.Error("list contacts", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, model.Contacts{Contacts: contacts})
}
