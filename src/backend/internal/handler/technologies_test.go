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

func TestGetTechnologies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mockFn     func(ctx context.Context) ([]model.Technology, error)
		wantStatus int
		wantCount  int
	}{
		{
			name: "success with multiple technologies",
			mockFn: func(ctx context.Context) ([]model.Technology, error) {
				return []model.Technology{
					{ID: 1, Name: "Go"},
					{ID: 2, Name: "React"},
					{ID: 3, Name: "PostgreSQL"},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  3,
		},
		{
			name: "success with empty list",
			mockFn: func(ctx context.Context) ([]model.Technology, error) {
				return []model.Technology{}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "store error returns 500",
			mockFn: func(ctx context.Context) ([]model.Technology, error) {
				return nil, errors.New("db unreachable")
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStore{ListTechnologiesFn: tc.mockFn}
			h := New(nil, store)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/technologies", nil)
			rec := httptest.NewRecorder()

			h.GetTechnologies(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantCount >= 0 {
				var got []model.Technology
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
