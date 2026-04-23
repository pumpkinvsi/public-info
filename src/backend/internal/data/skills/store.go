package skills

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/shared/db"
)

type Store interface {
	ListSkills(ctx context.Context) ([]Skill, error)
}

type repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) ListSkills(ctx context.Context) ([]Skill, error) {
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

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query skills: %w", err)
	}
	defer rows.Close()

	skills, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Skill, error) {
		var (
			skillName   string
			levelID     int
			levelWeight int
			levelTitle  string
		)
		if err := row.Scan(&skillName, &levelID, &levelWeight, &levelTitle); err != nil {
			return Skill{}, err
		}
		return Skill{
			Name: skillName,
			Level: Level{
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
