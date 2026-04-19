package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/model"
)

// ListSkills returns all skills joined with their proficiency level.
func (s *Store) ListSkills(ctx context.Context) ([]model.Skill, error) {
	const query = `
		SELECT
			s.name,
			l.id,
			l.level,
			l.title
		FROM skills s
		JOIN levels l ON l.id = s.level
		ORDER BY s.id
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query skills: %w", err)
	}
	defer rows.Close()

	skills, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Skill, error) {
		var (
			skillName   string
			levelID     int
			levelWeight int
			levelTitle  string
		)
		if err := row.Scan(&skillName, &levelID, &levelWeight, &levelTitle); err != nil {
			return model.Skill{}, err
		}
		return model.Skill{
			Name: skillName,
			Level: model.Level{
				ID:    levelID,
				Level: levelWeight,
				Text:  levelTitle,
			},
		}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect skills: %w", err)
	}

	return skills, nil
}
