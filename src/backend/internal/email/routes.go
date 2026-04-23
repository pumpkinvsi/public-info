package email

import (
	"src/backend/internal/shared/db"
	"src/backend/internal/shared/outbox"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, db *db.Postgres) {
	h := newHandler(outbox.NewRepository(db))

	r.Post("/email", h.SendEmail)
}
