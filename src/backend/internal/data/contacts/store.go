package contacts

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"src/backend/internal/shared/db"
)

type Store interface {
	ListContacts(ctx context.Context) ([]Contact, error)
}

type repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) ListContacts(ctx context.Context) ([]Contact, error) {
	const query = `SELECT name, value FROM contacts ORDER BY id`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query contacts: %w", err)
	}
	defer rows.Close()

	contacts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Contact, error) {
		var c Contact
		if err := row.Scan(&c.Name, &c.Value); err != nil {
			return Contact{}, err
		}
		return c, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect contacts: %w", err)
	}

	return contacts, nil
}
