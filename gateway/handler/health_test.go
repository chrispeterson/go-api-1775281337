package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthHandler uses a table-driven test pattern to cover multiple scenarios
func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
		description    string
	}{
		{
			name:           "GET request returns 200",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
			description:    "Should return 200 status code for GET /health",
		},
		{
			name:           "POST request returns 200",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
			description:    "Should return 200 status code for POST /health",
		},
		{
			name:           "HEAD request returns 200",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			description:    "Should return 200 status code for HEAD /health",
		},
	}

	for _, tt := range tests {
		t := tt // capture loop variable
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			// Call handler
			Health(rec, req)

			// Assert status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("%s: got status %d, want %d", tt.description, rec.Code, tt.expectedStatus)
			}

			// Assert Content-Type header for non-HEAD requests
			if tt.method != http.MethodHead {
				contentType := rec.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Content-Type: got %q, want "application/json"", contentType)
				}
			}

			// Assert response body for non-HEAD requests
			if tt.method != http.MethodHead && len(tt.expectedBody) > 0 {
				var response HealthResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
				}
				if response.Status != tt.expectedBody {
					t.Errorf("Status field: got %q, want %q", response.Status, tt.expectedBody)
				}
			}
		})
	}
}

// TestHealthHandlerStatusField verifies the status field is present and correct
func TestHealthHandlerStatusField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	var response HealthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Status field: got %q, want "ok"", response.Status)
	}
}

// TestHealthHandlerContentType verifies Content-Type header is set correctly
func TestHealthHandlerContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type: got %q, want "application/json"", contentType)
	}
}

// TestHealthHandlerJSONValidity validates that response is valid JSON
func TestHealthHandlerJSONValidity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response body must be valid JSON: %v", err)
	}

	// Verify structure
	if _, ok := response["status"]; !ok {
		t.Error("Response should contain 'status' key")
	}
}

// TestHealthHandlerStatusCode verifies 200 OK status
func TestHealthHandlerStatusCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code: got %d, want %d", rec.Code, http.StatusOK)
	}
}
