package projects

import (
	"context"
	"encoding/json"
	"fmt"

	"src/backend/internal/data/skills"
	"src/backend/internal/data/technologies"
	"src/backend/internal/shared/db"
	"src/backend/internal/shared/model"
)

type store interface {
	ListProjectsGrouped(ctx context.Context) ([]ProjectGroup, error)
}

type repository struct {
	db *db.Postgres
}

func newRepository(db *db.Postgres) *repository {
	return &repository{
		db: db,
	}
}

type projKey struct {
	techID    int
	projectID int
}

func (r *repository) ListProjectsGrouped(ctx context.Context) ([]ProjectGroup, error) {
	const query = `
		SELECT
			t.id          AS tech_id,
			t.name        AS tech_name,
			p.id          AS project_id,
			p.name        AS project_name,
			p.description AS project_desc,
			p.note        AS project_note,
			s.id          AS skill_id,
			s.name        AS skill_name,
			l.id          AS level_id,
			l.level       AS level_weight,
			l.title       AS level_title
		FROM technologies t
		JOIN project_groups  pg ON pg.technology_id = t.id
		JOIN projects         p ON p.id             = pg.project_id
		LEFT JOIN project_skills ps ON ps.project_id = p.id
		LEFT JOIN skills          s ON s.id           = ps.skill_id
		LEFT JOIN levels          l ON l.id           = s.level
		ORDER BY t.id, p.id, s.id
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query projects grouped: %w", err)
	}
	defer rows.Close()

	var (
		groups  []ProjectGroup
		techIdx = make(map[int]int)
		projIdx = make(map[projKey]int)
	)

	for rows.Next() {
		var (
			techID      int
			techName    string
			projectID   int
			nameRaw     []byte
			descRaw     []byte
			noteRaw     []byte
			skillID     *int
			skillName   *string
			levelID     *int
			levelWeight *int
			levelTitle  *string
		)

		if err := rows.Scan(
			&techID, &techName,
			&projectID, &nameRaw, &descRaw, &noteRaw,
			&skillID, &skillName,
			&levelID, &levelWeight, &levelTitle,
		); err != nil {
			return nil, fmt.Errorf("scan project row: %w", err)
		}

		tIdx, exists := techIdx[techID]
		if !exists {
			groups = append(groups, ProjectGroup{
				Technology: technologies.Technology{ID: techID, Name: techName},
				Projects:   []Project{},
			})
			tIdx = len(groups) - 1
			techIdx[techID] = tIdx
		}

		pk := projKey{techID: techID, projectID: projectID}
		pIdx, exists := projIdx[pk]
		if !exists {
			var name model.LocalizedString
			if err := json.Unmarshal(nameRaw, &name); err != nil {
				return nil, fmt.Errorf("unmarshal project name (id=%d): %w", projectID, err)
			}

			var desc model.LocalizedString
			if err := json.Unmarshal(descRaw, &desc); err != nil {
				return nil, fmt.Errorf("unmarshal project description (id=%d): %w", projectID, err)
			}

			var note *model.LocalizedString
			if noteRaw != nil {
				note = new(model.LocalizedString)
				if err := json.Unmarshal(noteRaw, note); err != nil {
					return nil, fmt.Errorf("unmarshal project note (id=%d): %w", projectID, err)
				}
			}

			groups[tIdx].Projects = append(groups[tIdx].Projects, Project{
				ID:          projectID,
				Name:        name,
				Description: desc,
				Note:        note,
				Skills:      []skills.Skill{},
			})
			pIdx = len(groups[tIdx].Projects) - 1
			projIdx[pk] = pIdx
		}

		if skillID != nil {
			groups[tIdx].Projects[pIdx].Skills = append(
				groups[tIdx].Projects[pIdx].Skills,
				skills.Skill{
					Name: *skillName,
					Level: skills.Level{
						ID:    *levelID,
						Level: *levelWeight,
						Text:  *levelTitle,
					},
				},
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project rows: %w", err)
	}

	return groups, nil
}
