package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"src/backend/internal/model"
)

func TestGetProjects(t *testing.T) {
	t.Parallel()

	note := model.LocalizedString{Rus: "Заметка", Eng: "Note"}

	tests := []struct {
		name       string
		mockFn     func(ctx context.Context) ([]model.ProjectGroup, error)
		wantStatus int
		wantCount  int
	}{
		{
			name: "success with groups and skills",
			mockFn: func(ctx context.Context) ([]model.ProjectGroup, error) {
				return []model.ProjectGroup{
					{
						Technology: model.Technology{ID: 1, Name: "Go"},
						Projects: []model.Project{
							{
								ID:          1,
								Name:        model.LocalizedString{Rus: "Проект", Eng: "Project"},
								Description: model.LocalizedString{Rus: "Описание", Eng: "Description"},
								Note:        &note,
								Skills: []model.Skill{
									{Name: "Go", Level: model.Level{ID: 1, Level: 3, Text: "Senior"}},
								},
							},
						},
					},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "success with nil note",
			mockFn: func(ctx context.Context) ([]model.ProjectGroup, error) {
				return []model.ProjectGroup{
					{
						Technology: model.Technology{ID: 1, Name: "Go"},
						Projects: []model.Project{
							{
								ID:          2,
								Name:        model.LocalizedString{Eng: "No note project"},
								Description: model.LocalizedString{Eng: "Desc"},
								Note:        nil,
								Skills:      []model.Skill{},
							},
						},
					},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "success with empty result",
			mockFn: func(ctx context.Context) ([]model.ProjectGroup, error) {
				return []model.ProjectGroup{}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "store error returns 500",
			mockFn: func(ctx context.Context) ([]model.ProjectGroup, error) {
				return nil, errors.New("query failed")
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStore{ListProjectsGroupedFn: tc.mockFn}
			h := New(nil, store)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
			rec := httptest.NewRecorder()

			h.GetProjects(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantCount >= 0 {
				var got []model.ProjectGroup
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(got) != tc.wantCount {
					t.Errorf("groups count: got %d, want %d", len(got), tc.wantCount)
				}
			}
		})
	}
}

func TestGetProjectsNilNoteSerialisation(t *testing.T) {
	t.Parallel()

	store := &mockStore{
		ListProjectsGroupedFn: func(ctx context.Context) ([]model.ProjectGroup, error) {
			return []model.ProjectGroup{
				{
					Technology: model.Technology{ID: 1, Name: "Go"},
					Projects: []model.Project{
						{ID: 1, Note: nil, Skills: []model.Skill{}},
					},
				},
			}, nil
		},
	}
	h := New(nil, store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	h.GetProjects(rec, req)

	var raw []struct {
		Projects []struct {
			Note any `json:"note"`
		} `json:"projects"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(raw) == 0 || len(raw[0].Projects) == 0 {
		t.Fatal("expected at least one group with one project")
	}
	if raw[0].Projects[0].Note != nil {
		t.Errorf("note: got %v, want null", raw[0].Projects[0].Note)
	}
}
