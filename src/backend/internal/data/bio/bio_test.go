package bio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"

	"src/backend/internal/shared/model"
)

type mockStore struct {
	GetBioFn func(ctx context.Context) (*Bio, error)
}

func (m *mockStore) GetBio(ctx context.Context) (*Bio, error) {
	return m.GetBioFn(ctx)
}

func TestGetBio(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mockFn     func(ctx context.Context) (*Bio, error)
		wantStatus int
		wantBio    *Bio
	}{
		{
			name: "success",
			mockFn: func(ctx context.Context) (*Bio, error) {
				return &Bio{
					Text: model.LocalizedString{
						Rus: "Обо мне",
						Eng: "About me",
					},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantBio: &Bio{
				Text: model.LocalizedString{
					Rus: "Обо мне",
					Eng: "About me",
				},
			},
		},
		{
			name: "not found returns 404",
			mockFn: func(ctx context.Context) (*Bio, error) {
				return nil, pgx.ErrNoRows
			},
			wantStatus: http.StatusNotFound,
			wantBio:    nil,
		},
		{
			name: "store error returns 500",
			mockFn: func(ctx context.Context) (*Bio, error) {
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
			h := newHandler(store)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/bio", nil)
			rec := httptest.NewRecorder()

			h.GetBio(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantBio != nil {
				var got Bio
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
