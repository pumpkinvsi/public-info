package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// GET /api/v1/bio
func (h *Handler) GetBio(w http.ResponseWriter, r *http.Request) {
	bio, err := h.store.GetBio(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "bio not found")
			return
		}
		slog.Error("get bio", "error", err)
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, http.StatusOK, bio)
}
