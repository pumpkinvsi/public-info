package handler

import (
	"encoding/json"
	"net/http"

	"src/backend/internal/model"
)

// POST /api/v1/email
func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var payload model.Email

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// TODO: validate
	// TODO: publish to Kafka
	w.WriteHeader(http.StatusAccepted)
}