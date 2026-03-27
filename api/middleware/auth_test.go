package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"okapi/internal/config"
)

func TestAuth(t *testing.T) {
	tests := []struct {
		name           string
		cfg            config.AuthConfig
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:           "Auth disabled",
			cfg:            config.AuthConfig{Enabled: false},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Auth enabled, no key",
			cfg:            config.AuthConfig{Enabled: true, APIKeys: []string{"secret"}},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Auth enabled, valid X-API-Key",
			cfg:            config.AuthConfig{Enabled: true, APIKeys: []string{"secret"}},
			headers:        map[string]string{"X-API-Key": "secret"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Auth enabled, valid Bearer token",
			cfg:            config.AuthConfig{Enabled: true, APIKeys: []string{"secret"}},
			headers:        map[string]string{"Authorization": "Bearer secret"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Auth enabled, invalid key",
			cfg:            config.AuthConfig{Enabled: true, APIKeys: []string{"secret"}},
			headers:        map[string]string{"X-API-Key": "wrong"},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Auth(&tt.cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
