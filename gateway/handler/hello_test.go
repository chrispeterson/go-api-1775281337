package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrispeterson/go-api-1775281337/gateway/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloHandler_Unit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		url        string
		method     string
		wantStatus int
		wantJSON   map[string]string
	}{
		{
			name:       "name=Alice returns 200 greeting",
			url:        "/hello?name=Alice",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantJSON:   map[string]string{"message": "Hello, Alice!"},
		},
		{
			name:       "name=Bob returns 200 greeting",
			url:        "/hello?name=Bob",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantJSON:   map[string]string{"message": "Hello, Bob!"},
		},
		{
			name:       "name with special chars returns 200 greeting",
			url:        "/hello?name=O%27Brien",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantJSON:   map[string]string{"message": "Hello, O'Brien!"},
		},
		{
			name:       "empty name param returns 400",
			url:        "/hello?name=",
			method:     http.MethodGet,
			wantStatus: http.StatusBadRequest,
			wantJSON:   map[string]string{"error": "name parameter is required"},
		},
		{
			name:       "missing name param returns 400",
			url:        "/hello",
			method:     http.MethodGet,
			wantStatus: http.StatusBadRequest,
			wantJSON:   map[string]string{"error": "name parameter is required"},
		},
		{
			name:       "POST method returns 405",
			url:        "/hello?name=Alice",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
			wantJSON:   map[string]string{"error": "method not allowed"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := handler.NewHelloHandler()
			req := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.wantStatus, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			var got map[string]string
			require.NoError(t, json.Unmarshal(body, &got))
			assert.Equal(t, tc.wantJSON, got)
		})
	}
}

func TestHelloHandler_Integration(t *testing.T) {
	t.Parallel()

	h := handler.NewHelloHandler()
	svr := httptest.NewServer(h)
	defer svr.Close()

	baseURL := svr.URL

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantKey    string
		wantValue  string
	}{
		{
			name:       "Alice gets greeted",
			path:       "/?name=Alice",
			wantStatus: http.StatusOK,
			wantKey:    "message",
			wantValue:  "Hello, Alice!",
		},
		{
			name:       "missing param returns 400",
			path:       "/",
			wantStatus: http.StatusBadRequest,
			wantKey:    "error",
			wantValue:  "name parameter is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := http.Get(baseURL + tc.path)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tc.wantStatus, res.StatusCode)

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			var got map[string]string
			require.NoError(t, json.Unmarshal(body, &got))
			assert.Equal(t, tc.wantValue, got[tc.wantKey])
		})
	}
}
