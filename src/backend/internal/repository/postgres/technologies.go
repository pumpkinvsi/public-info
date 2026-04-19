package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/model"
)

// ListTechnologies returns all technology entries ordered by id.
func (s *Store) ListTechnologies(ctx context.Context) ([]model.Technology, error) {
	const query = `SELECT id, name FROM technologies ORDER BY id`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query technologies: %w", err)
	}
	defer rows.Close()

	techs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Technology, error) {
		var t model.Technology
		if err := row.Scan(&t.ID, &t.Name); err != nil {
			return model.Technology{}, err
		}
		return t, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect technologies: %w", err)
	}

	return techs, nil
}
