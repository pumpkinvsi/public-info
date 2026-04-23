package contacts

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStore struct {
	ListContactsFn func(ctx context.Context) ([]Contact, error)
}

func (m *mockStore) ListContacts(ctx context.Context) ([]Contact, error) {
	return m.ListContactsFn(ctx)
}

func TestGetContacts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mockFn     func(ctx context.Context) ([]Contact, error)
		wantStatus int
		wantCount  int
	}{
		{
			name: "success with multiple contacts",
			mockFn: func(ctx context.Context) ([]Contact, error) {
				return []Contact{
					{Name: "Email", Value: "me@example.com"},
					{Name: "GitHub", Value: "github.com/me"},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "success with empty list",
			mockFn: func(ctx context.Context) ([]Contact, error) {
				return []Contact{}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "store error returns 500",
			mockFn: func(ctx context.Context) ([]Contact, error) {
				return nil, errors.New("scan error")
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStore{ListContactsFn: tc.mockFn}
			h := NewHandler(store)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/contacts", nil)
			rec := httptest.NewRecorder()

			h.GetContacts(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantCount >= 0 {
				var got Contacts
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(got.Contacts) != tc.wantCount {
					t.Errorf("count: got %d, want %d", len(got.Contacts), tc.wantCount)
				}
			}
		})
	}
}

func TestGetContactsResponseEnvelope(t *testing.T) {
	t.Parallel()

	store := &mockStore{
		ListContactsFn: func(ctx context.Context) ([]Contact, error) {
			return []Contact{{Name: "Email", Value: "me@example.com"}}, nil
		},
	}
	h := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/contacts", nil)
	rec := httptest.NewRecorder()
	h.GetContacts(rec, req)

	var raw map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := raw["contacts"]; !ok {
		t.Error("response envelope missing 'contacts' key")
	}
}
