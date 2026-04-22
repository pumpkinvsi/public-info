package repository

import (
	"context"

	"src/backend/internal/model"
)

type Health interface {
	Ping(ctx context.Context) error
}

type Bio interface {
	GetBio(ctx context.Context) (*model.Bio, error)
}

type Skills interface {
	ListSkills(ctx context.Context) ([]model.Skill, error)
}

type Projects interface {
	ListProjectsGrouped(ctx context.Context) ([]model.ProjectGroup, error)
}

type Technologies interface {
	ListTechnologies(ctx context.Context) ([]model.Technology, error)
}

type Contacts interface {
	ListContacts(ctx context.Context) ([]model.Contact, error)
}

type Store interface {
	Health
	Bio
	Skills
	Projects
	Technologies
	Contacts
}
