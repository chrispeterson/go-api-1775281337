package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth_GetRequest_Returns200(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code should be 200")
}

func TestHealth_ResponseBody_ValidJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err, "response body should be valid JSON")
	assert.Equal(t, "ok", responseBody["status"], "status field should be 'ok'")
}

func TestHealth_ContentTypeHeader_ApplicationJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", contentType, "Content-Type header should be application/json")
}

func TestHealth_TableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectBody     bool
	}{
		{
			name:           "GET request returns 200 with body",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "POST request returns 405",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectBody:     false,
		},
		{
			name:           "PUT request returns 405",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectBody:     false,
		},
		{
			name:           "DELETE request returns 405",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
			expectBody:     false,
		},
		{
			name:           "PATCH request returns 405",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			Health(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "status code mismatch")

			if tt.expectBody {
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err, "response should be valid JSON")
				assert.Equal(t, "ok", responseBody["status"])
			}
		})
	}
}

func TestHealth_ConsistentAcrossMultipleCalls(t *testing.T) {
	t.Parallel()

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		Health(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "call %d: status should be 200", i+1)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "call %d: Content-Type mismatch", i+1)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err, "call %d: response should be valid JSON", i+1)
		assert.Equal(t, "ok", responseBody["status"], "call %d: status field mismatch", i+1)
	}
}

func TestHealth_ResponseBodyStructure(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err, "response should be valid JSON")

	// Verify required fields exist
	status, exists := responseBody["status"]
	require.True(t, exists, "response must contain 'status' field")
	assert.Equal(t, "ok", status, "status field should be 'ok'")

	// Verify no extra fields
	assert.Equal(t, 1, len(responseBody), "response should contain exactly one field")
}

func TestHealth_HeadersSetBeforeWriteHeader(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	// Verify Content-Type is present
	headers := w.Header()
	assert.NotNil(t, headers)
	assert.NotEmpty(t, headers.Get("Content-Type"))
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
}

func TestHealth_MethodNotAllowed_AllMethods(t *testing.T) {
	t.Parallel()

	notAllowedMethods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range notAllowedMethods {
		t.Run("method "+method, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(method, "/health", nil)
			w := httptest.NewRecorder()

			Health(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code, "method %s should not be allowed", method)
		})
	}
}
