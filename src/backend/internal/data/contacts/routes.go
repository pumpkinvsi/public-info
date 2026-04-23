package contacts

import (
	"github.com/go-chi/chi/v5"
	
	"src/backend/internal/shared/db"
)

func RegisterRoutes(r chi.Router,db *db.Postgres) {
	h := NewHandler(NewRepository(db))

	r.Get("/contacts", h.GetContacts)
}
