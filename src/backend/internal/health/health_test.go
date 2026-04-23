package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStore struct {
	PingFn func(ctx context.Context) error
}

func (m *mockStore) Ping(ctx context.Context) error {
	return m.PingFn(ctx)
}

func TestLiveness(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rec := httptest.NewRecorder()

	h.Liveness(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	var got healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Status != "ok" {
		t.Errorf("status field: got %q, want %q", got.Status, "ok")
	}
}

func TestReadiness(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pingFn         func(ctx context.Context) error
		wantStatus     int
		wantOverall    string
		wantDBStatus   string
		wantDBErrEmpty bool
	}{
		{
			name:           "database reachable returns 200",
			pingFn:         func(ctx context.Context) error { return nil },
			wantStatus:     http.StatusOK,
			wantOverall:    "ok",
			wantDBStatus:   "ok",
			wantDBErrEmpty: true,
		},
		{
			name:           "database unreachable returns 503",
			pingFn:         func(ctx context.Context) error { return errors.New("connection refused") },
			wantStatus:     http.StatusServiceUnavailable,
			wantOverall:    "degraded",
			wantDBStatus:   "degraded",
			wantDBErrEmpty: false,
		},
		{
			name:           "context cancelled returns 503",
			pingFn:         func(ctx context.Context) error { return context.Canceled },
			wantStatus:     http.StatusServiceUnavailable,
			wantOverall:    "degraded",
			wantDBStatus:   "degraded",
			wantDBErrEmpty: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &mockStore{PingFn: tc.pingFn}
			h := NewHandler(store)

			req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
			rec := httptest.NewRecorder()

			h.Readiness(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}

			var got healthResponse
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decode: %v", err)
			}

			if got.Status != tc.wantOverall {
				t.Errorf("overall status: got %q, want %q", got.Status, tc.wantOverall)
			}

			db, ok := got.Checks["database"]
			if !ok {
				t.Fatal("response missing 'database' check")
			}

			if db.Status != tc.wantDBStatus {
				t.Errorf("db status: got %q, want %q", db.Status, tc.wantDBStatus)
			}

			if tc.wantDBErrEmpty && db.Error != "" {
				t.Errorf("db error: got %q, want empty", db.Error)
			}
			if !tc.wantDBErrEmpty && db.Error == "" {
				t.Error("db error: got empty, want non-empty")
			}
		})
	}
}
