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

func TestGetSkills(t *testing.T) {
	t.Parallel()

	senior := model.Level{ID: 1, Level: 3, Text: "Senior"}

	tests := []struct {
		name       string
		mockFn     func(ctx context.Context) ([]model.Skill, error)
		wantStatus int
		wantCount  int
	}{
		{
			name: "success with multiple skills",
			mockFn: func(ctx context.Context) ([]model.Skill, error) {
				return []model.Skill{
					{Name: "Go", Level: senior},
					{Name: "PostgreSQL", Level: senior},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "success with empty list",
			mockFn: func(ctx context.Context) ([]model.Skill, error) {
				return []model.Skill{}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "store error returns 500",
			mockFn: func(ctx context.Context) ([]model.Skill, error) {
				return nil, errors.New("timeout")
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStore{ListSkillsFn: tc.mockFn}
			h := New(nil, store)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/skills", nil)
			rec := httptest.NewRecorder()

			h.GetSkills(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantCount >= 0 {
				var got []model.Skill
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(got) != tc.wantCount {
					t.Errorf("count: got %d, want %d", len(got), tc.wantCount)
				}
			}
		})
	}
}
