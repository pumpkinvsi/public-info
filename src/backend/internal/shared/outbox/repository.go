package outbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/shared/db"
)

type Outbox interface {
	Insert(ctx context.Context, payload []byte) error
	ProcessNext(ctx context.Context, handler func(msg *Message) error) (bool, error)
}

type Message struct {
	ID      int
	Payload []byte
}

type repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Insert(ctx context.Context, payload []byte) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO outbox (event_type, payload) VALUES ($1, $2)`,
		"EmailReceived", payload,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) ProcessNext(ctx context.Context, handler func(msg *Message) error) (bool, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `SELECT id, payload FROM outbox WHERE processed = false LIMIT 1`)
	var msg Message
	if err := row.Scan(&msg.ID, &msg.Payload); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if err := handler(&msg); err != nil {
		return false, err
	}

	_, err = tx.Exec(ctx, `UPDATE outbox SET processed = true WHERE id = $1`, msg.ID)
	if err != nil {
		return false, err
	}

	return true, tx.Commit(ctx)
}
