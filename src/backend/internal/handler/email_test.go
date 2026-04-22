package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "success with valid payload",
			body:       `{"text":"Hello","sender":"Alice","contact":"alice@example.com"}`,
			wantStatus: http.StatusAccepted,
		},
		{
			name:       "success ignores unknown fields",
			body:       `{"text":"Hi","sender":"Bob","contact":"bob@example.com","extra":"ignored"}`,
			wantStatus: http.StatusAccepted,
		},
		{
			name:       "malformed JSON returns 400",
			body:       `{not valid json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body returns 400",
			body:       ``,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong JSON type returns 400",
			body:       `["text","sender"]`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := New(nil, &mockStore{})

			req := httptest.NewRequest(http.MethodPost, "/api/v1/email", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.SendEmail(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}
