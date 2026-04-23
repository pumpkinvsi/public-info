package skills

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

// GET /api/v1/skills
func (h *handler) GetSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := h.store.ListSkills(r.Context())
	if err != nil {
		slog.Error("list skills", "error", err)
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpUtils.RespondJSON(w, http.StatusOK, skills)
}
