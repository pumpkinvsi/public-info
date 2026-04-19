package handler

import (
	"net/http"
)

type healthResponse struct {
	Status string                 `json:"status"`
	Checks map[string]checkResult `json:"checks,omitempty"`
}

type checkResult struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// GET /health/live
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

// GET /health/ready
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	db := checkResult{Status: "ok"} // TODO: replace with real db check

	status := http.StatusOK
	overall := "ok"

	if db.Status != "ok" {
		overall = "degraded"
		status = http.StatusServiceUnavailable
	}

	respondJSON(w, status, healthResponse{
		Status: overall,
		Checks: map[string]checkResult{
			"database": db,
		},
	})
}