package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/model"
)

func TestGetBio(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mockFn     func(ctx context.Context) (*model.Bio, error)
		wantStatus int
		wantBio    *model.Bio
	}{
		{
			name: "success",
			mockFn: func(ctx context.Context) (*model.Bio, error) {
				return &model.Bio{
					Text: model.LocalizedString{
						Rus: "Обо мне",
						Eng: "About me",
					},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantBio: &model.Bio{
				Text: model.LocalizedString{
					Rus: "Обо мне",
					Eng: "About me",
				},
			},
		},
		{
			name: "not found returns 404",
			mockFn: func(ctx context.Context) (*model.Bio, error) {
				return nil, pgx.ErrNoRows
			},
			wantStatus: http.StatusNotFound,
			wantBio:    nil,
		},
		{
			name: "store error returns 500",
			mockFn: func(ctx context.Context) (*model.Bio, error) {
				return nil, errors.New("connection reset")
			},
			wantStatus: http.StatusInternalServerError,
			wantBio:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStore{GetBioFn: tc.mockFn}
			h := New(nil, store)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/bio", nil)
			rec := httptest.NewRecorder()

			h.GetBio(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantBio != nil {
				var got model.Bio
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if got.Text.Eng != tc.wantBio.Text.Eng || got.Text.Rus != tc.wantBio.Text.Rus {
					t.Errorf("body: got %+v, want %+v", got, tc.wantBio)
				}
			}
		})
	}
}
