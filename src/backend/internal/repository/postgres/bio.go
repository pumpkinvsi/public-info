package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"src/backend/internal/model"
)

func (s *Store) GetBio(ctx context.Context) (*model.Bio, error) {
	const query = `SELECT bio FROM info WHERE id = 1`

	var raw []byte
	if err := s.pool.QueryRow(ctx, query).Scan(&raw); err != nil {
		return nil, fmt.Errorf("query bio: %w", err)
	}

	var ls model.LocalizedString
	if err := json.Unmarshal(raw, &ls); err != nil {
		return nil, fmt.Errorf("unmarshal bio: %w", err)
	}

	return &model.Bio{Text: ls}, nil
}
