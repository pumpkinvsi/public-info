package email

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"src/backend/internal/shared/outbox"
	httpUtils "src/backend/internal/shared/http"
)

type handler struct {
	outbox outbox.Outbox
}

func NewHandler(outbox outbox.Outbox) *handler {
	return &handler{
		outbox: outbox,
	}
}

const maxTextLength = 5000

// POST /api/v1/email
func (h *handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var payload Email

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpUtils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := validateEmail(payload); err != nil {
		httpUtils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.outbox.Insert(r.Context(), raw); err != nil {
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func validateEmail(e Email) error {
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
