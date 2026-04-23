package technologies

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/shared/db"
)

type Store interface {
	ListTechnologies(ctx context.Context) ([]Technology, error)
}

type repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) ListTechnologies(ctx context.Context) ([]Technology, error) {
	const query = `SELECT id, name FROM technologies ORDER BY id`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query technologies: %w", err)
	}
	defer rows.Close()

	techs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Technology, error) {
		var t Technology
		if err := row.Scan(&t.ID, &t.Name); err != nil {
			return Technology{}, err
		}
		return t, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect technologies: %w", err)
	}

	return techs, nil
}
