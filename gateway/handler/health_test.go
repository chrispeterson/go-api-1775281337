package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name           string
		wantStatus     int
		wantStatusBody string
	}{
		{
			name:           "successful request returns 200",
			wantStatus:     http.StatusOK,
			wantStatusBody: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/health", nil)

			Health(rr, req)

			// Check status code
			if rr.Code != tt.wantStatus {
				t.Errorf("Health() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			// Check Content-Type header
			if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Health() Content-Type = %v, want application/json", ct)
			}

			// Check response body is valid JSON
			var resp HealthResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Errorf("Health() response body is not valid JSON: %v", err)
			}

			// Check response contains correct status field
			if resp.Status != tt.wantStatusBody {
				t.Errorf("Health() status field = %v, want %v", resp.Status, tt.wantStatusBody)
			}
		})
	}
}
