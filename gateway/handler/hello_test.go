package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelloResponse represents the expected JSON response structure
type TestHelloResponse struct {
	Message string `json:"message"`
}

// TestErrorResponse represents the error response structure
type TestErrorResponse struct {
	Error string `json:"error"`
}

// TestHelloEndpoint is a table-driven test for the /hello endpoint
func TestHelloEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		expectedMsg    string
		isError        bool
		description    string
	}{
		{
			name:           "hello with name=Alice",
			queryParam:     "name=Alice",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, Alice!",
			isError:        false,
			description:    "GET /hello?name=Alice returns personalized greeting",
		},
		{
			name:           "hello without name parameter",
			queryParam:     "",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, World!",
			isError:        false,
			description:    "GET /hello (no parameter) returns default greeting",
		},
		{
			name:           "hello with empty name parameter",
			queryParam:     "name=",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, World!",
			isError:        false,
			description:    "GET /hello?name= (empty parameter) returns default greeting",
		},
		{
			name:           "hello with name=Bob",
			queryParam:     "name=Bob",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, Bob!",
			isError:        false,
			description:    "GET /hello?name=Bob returns correct greeting",
		},
	}

	for _, tt := range tests {
		tt := tt // capture for parallel execution
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			var req *http.Request
			if tt.queryParam != "" {
				req = httptest.NewRequest("GET", "/hello?"+tt.queryParam, nil)
			} else {
				req = httptest.NewRequest("GET", "/hello", nil)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create handler and call it
			handler := NewHelloHandler()
			handler.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tt.expectedStatus, w.Code, tt.description)

			// Verify Content-Type header
			contentType := w.Header().Get("Content-Type")
			assert.Equal(t, "application/json", contentType,
				"response should have application/json content type")

			// Parse and verify response
			var resp TestHelloResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err, "response should be valid JSON")

			// Verify message content
			assert.Equal(t, tt.expectedMsg, resp.Message,
				"response message should match expected greeting")
		})
	}
}

// TestHelloResponseStructure validates JSON response structure and format
func TestHelloResponseStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/hello?name=TestUser", nil)
	w := httptest.NewRecorder()

	handler := NewHelloHandler()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "handler should return 200 OK")

	var resp TestHelloResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "response should be valid JSON")

	require.NotEmpty(t, resp.Message, "message field should not be empty")
	assert.Contains(t, resp.Message, "Hello", "message should contain greeting")
	assert.Contains(t, resp.Message, "TestUser", "message should contain the name")
}

// TestHelloContentTypeHeader verifies correct Content-Type header
func TestHelloContentTypeHeader(t *testing.T) {
	tests := []struct {
		name      string
		queryStr  string
	}{
		{"with name", "?name=Alice"},
		{"without name", ""},
		{"empty name", "?name="},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello"+tt.queryStr, nil)
			w := httptest.NewRecorder()

			handler := NewHelloHandler()
			handler.ServeHTTP(w, req)

			contentType := w.Header().Get("Content-Type")
			assert.Equal(t, "application/json", contentType,
				"content type should always be application/json")
		})
	}
}

// TestHelloStatusCode200 verifies HTTP 200 status code for valid requests
func TestHelloStatusCode200(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"with name=Alice", "/hello?name=Alice"},
		{"with name=Bob", "/hello?name=Bob"},
		{"with name=TestUser", "/hello?name=TestUser"},
		{"without name (default to World)", "/hello"},
		{"with empty name parameter", "/hello?name="},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			handler := NewHelloHandler()
			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code,
				"handler should always return 200 OK for valid requests")
		})
	}
}

// TestHelloWithNameParameter tests the name parameter parsing
func TestHelloWithNameParameter(t *testing.T) {
	tests := []struct {
		name          string
		nameParam     string
		expectedGreet string
	}{
		{"Alice", "Alice", "Hello, Alice!"},
		{"Bob", "Bob", "Hello, Bob!"},
		{"Charlie", "Charlie", "Hello, Charlie!"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello?name="+tt.nameParam, nil)
			w := httptest.NewRecorder()

			handler := NewHelloHandler()
			handler.ServeHTTP(w, req)

			var resp TestHelloResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedGreet, resp.Message,
				"message should contain the correct name")
		})
	}
}

// TestHelloDefaultWorld tests that missing name defaults to World
func TestHelloDefaultWorld(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		description string
	}{
		{"no query string", "/hello", "should use World when no name provided"},
		{"empty name", "/hello?name=", "should use World when name is empty"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			handler := NewHelloHandler()
			handler.ServeHTTP(w, req)

			var resp TestHelloResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			assert.Equal(t, "Hello, World!", resp.Message, tt.description)
		})
	}
}
