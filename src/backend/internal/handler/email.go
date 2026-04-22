package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"src/backend/internal/model"
)

const maxTextLength = 5000

// POST /api/v1/email
func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var payload model.Email

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := validateEmail(payload); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: publish to Kafka
	w.WriteHeader(http.StatusAccepted)
}

func validateEmail(e model.Email) error {
	if strings.TrimSpace(e.Sender) == "" {
		return errors.New("sender is required")
	}
	if strings.TrimSpace(e.Contact) == "" {
		return errors.New("contact is required")
	}
	if _, err := mail.ParseAddress(e.Contact); err != nil {
		return errors.New("contact must be a valid email address")
	}
	if strings.TrimSpace(e.Text) == "" {
		return errors.New("text is required")
	}
	if len([]rune(e.Text)) > maxTextLength {
		return errors.New("text exceeds maximum allowed length")
	}
	return nil
}
