package health

import (
	"github.com/go-chi/chi/v5"

	"src/backend/internal/shared/db"
)

func RegisterRoutes(r chi.Router, db *db.Postgres) {
	h := NewHandler(NewRepository(db))
	
	r.Route("/health", func(r chi.Router) {
		r.Get("/live", h.Liveness)
		r.Get("/ready", h.Readiness)
	})
}