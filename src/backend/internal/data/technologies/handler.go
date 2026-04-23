package technologies

import (
	"log/slog"
	"net/http"

	httpUtils "src/backend/internal/shared/http"
)

type handler struct {
	store Store
}

func NewHandler(store Store) *handler {
	return &handler{
		store: store,
	}
}

// GET /api/v1/technologies
func (h *handler) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	techs, err := h.store.ListTechnologies(r.Context())
	if err != nil {
		slog.Error("list technologies", "error", err)
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpUtils.RespondJSON(w, http.StatusOK, techs)
}
