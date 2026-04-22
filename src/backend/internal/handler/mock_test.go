package handler

import (
	"context"

	"src/backend/internal/model"
)

type mockStore struct {
	PingFn                func(ctx context.Context) error
	GetBioFn              func(ctx context.Context) (*model.Bio, error)
	ListSkillsFn          func(ctx context.Context) ([]model.Skill, error)
	ListProjectsGroupedFn func(ctx context.Context) ([]model.ProjectGroup, error)
	ListTechnologiesFn    func(ctx context.Context) ([]model.Technology, error)
	ListContactsFn        func(ctx context.Context) ([]model.Contact, error)
}

func (m *mockStore) Ping(ctx context.Context) error {
	return m.PingFn(ctx)
}

func (m *mockStore) GetBio(ctx context.Context) (*model.Bio, error) {
	return m.GetBioFn(ctx)
}

func (m *mockStore) ListSkills(ctx context.Context) ([]model.Skill, error) {
	return m.ListSkillsFn(ctx)
}

func (m *mockStore) ListProjectsGrouped(ctx context.Context) ([]model.ProjectGroup, error) {
	return m.ListProjectsGroupedFn(ctx)
}

func (m *mockStore) ListTechnologies(ctx context.Context) ([]model.Technology, error) {
	return m.ListTechnologiesFn(ctx)
}

func (m *mockStore) ListContacts(ctx context.Context) ([]model.Contact, error) {
	return m.ListContactsFn(ctx)
}
