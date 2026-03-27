package adapters

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"okapi/internal/models"
)

func TestGenericHTTPAdapter_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus models.Status
	}{
		{
			name:           "Operational status",
			statusCode:     http.StatusOK,
			expectedStatus: models.StatusOperational,
		},
		{
			name:           "Major outage status (500)",
			statusCode:     http.StatusInternalServerError,
			expectedStatus: models.StatusMajorOutage,
		},
		{
			name:           "Major outage status (404)",
			statusCode:     http.StatusNotFound,
			expectedStatus: models.StatusMajorOutage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := StatuspageConfig{
				ID:                  "test-http",
				DisplayName:         "Test HTTP",
				Kind:                "http",
				Subdomain:           server.URL,
				PollIntervalSeconds: 60,
			}
			adapter := NewGenericHTTPAdapter(cfg)

			res, err := adapter.Fetch(context.Background())
			if err != nil {
				t.Fatalf("Fetch failed: %v", err)
			}

			if res.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, res.Status)
			}

			if res.Service != cfg.ID {
				t.Errorf("expected service %s, got %s", cfg.ID, res.Service)
			}
		})
	}
}

func TestGenericHTTPAdapter_Fetch_Error(t *testing.T) {
	cfg := StatuspageConfig{
		ID:                  "test-http-error",
		Subdomain:           "invalid-url",
		PollIntervalSeconds: 60,
	}
	adapter := NewGenericHTTPAdapter(cfg)

	_, err := adapter.Fetch(context.Background())
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}
