package repository

import (
	"context"

	"src/backend/internal/model"
)

// Health is implemented by any data store that can report its own availability.
type Health interface {
	Ping(ctx context.Context) error
}

// Bio provides access to the singleton biography record.
type Bio interface {
	GetBio(ctx context.Context) (*model.Bio, error)
}

// Skills provides access to the skills catalog.
type Skills interface {
	ListSkills(ctx context.Context) ([]model.Skill, error)
}

// Projects provides access to portfolio projects.
type Projects interface {
	ListProjectsGrouped(ctx context.Context) ([]model.ProjectGroup, error)
}

// Technologies provides access to the technologies catalog.
type Technologies interface {
	ListTechnologies(ctx context.Context) ([]model.Technology, error)
}

// Contacts provides access to contact entries.
type Contacts interface {
	ListContacts(ctx context.Context) ([]model.Contact, error)
}

// Store is the aggregate repository interface used for dependency injection.
// Any type that satisfies all sub-interfaces can be passed as a Store.
type Store interface {
	Health
	Bio
	Skills
	Projects
	Technologies
	Contacts
}
