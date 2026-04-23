package projects

import (
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

// GET /api/v1/projects
func (h *handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.ListProjectsGrouped(r.Context())
	if err != nil {
		httpUtils.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpUtils.RespondJSON(w, http.StatusOK, projects)
}
