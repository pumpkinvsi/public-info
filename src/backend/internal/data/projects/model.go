package projects

import (
	"src/backend/internal/data/skills"
	"src/backend/internal/data/technologies"
	"src/backend/internal/shared/model"
)

type Project struct {
	ID          int              `json:"id"`
	Name        model.LocalizedString  `json:"name"`
	Description model.LocalizedString  `json:"description"`
	Skills      []skills.Skill          `json:"skills"`
	Note        *model.LocalizedString `json:"note"`
}

type ProjectGroup struct {
	Technology technologies.Technology `json:"technology"`
	Projects   []Project  `json:"projects"`
}
