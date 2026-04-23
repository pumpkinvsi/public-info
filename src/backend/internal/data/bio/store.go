package bio

import (
	"context"
	"encoding/json"
	"fmt"

	"src/backend/internal/shared/model"
	"src/backend/internal/shared/db"
)

type Store interface {
	GetBio(ctx context.Context) (*Bio, error)
}

type repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetBio(ctx context.Context) (*Bio, error) {
	const query = `SELECT bio FROM info WHERE id = 1`

	var raw []byte
	if err := r.db.Pool.QueryRow(ctx, query).Scan(&raw); err != nil {
		return nil, fmt.Errorf("query bio: %w", err)
	}

	var ls model.LocalizedString
	if err := json.Unmarshal(raw, &ls); err != nil {
		return nil, fmt.Errorf("unmarshal bio: %w", err)
	}

	return &Bio{Text: ls}, nil
}
