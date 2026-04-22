package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/model"
)

func (s *Store) ListContacts(ctx context.Context) ([]model.Contact, error) {
	const query = `SELECT name, value FROM contacts ORDER BY id`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query contacts: %w", err)
	}
	defer rows.Close()

	contacts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Contact, error) {
		var c model.Contact
		if err := row.Scan(&c.Name, &c.Value); err != nil {
			return model.Contact{}, err
		}
		return c, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect contacts: %w", err)
	}

	return contacts, nil
}
