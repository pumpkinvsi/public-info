package projects

import (
	"github.com/go-chi/chi/v5"

	"src/backend/internal/shared/db"
)

func RegisterRoutes(r chi.Router, db *db.Postgres) {
	h := newHandler(newRepository(db))

	r.Get("/projects", h.GetProjects)
}
