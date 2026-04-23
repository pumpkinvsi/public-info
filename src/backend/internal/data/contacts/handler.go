package contacts

import (
	"log/slog"
	"net/http"

	httpUtils "src/backend/internal/shared/http"
)

type handler struct {
	store store
}

func newHandler(store store) *handler {
	return &handler{
		store: store,
	}
}

// GET /api/v1/contacts
func (h *handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.store.ListContacts(r.Context())
	if err != nil {
		slog.Error("list contacts", "error", err)
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpUtils.RespondJSON(w, http.StatusOK, Contacts{Contacts: contacts})
}
