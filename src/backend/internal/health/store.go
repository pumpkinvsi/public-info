package health

import (
	"context"

	"src/backend/internal/shared/db"
)

type Store interface {
	Ping(ctx context.Context) error
}

type Repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.Pool.Ping(ctx)
}
