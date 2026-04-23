package bio

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"

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

// GET /api/v1/bio
func (h *handler) GetBio(w http.ResponseWriter, r *http.Request) {
	bio, err := h.store.GetBio(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpUtils.RespondError(w, http.StatusNotFound, "bio not found")
			return
		}
		slog.Error("get bio", "error", err)
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpUtils.RespondJSON(w, http.StatusOK, bio)
}
