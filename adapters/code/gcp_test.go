package code

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"okapi/internal/models"
)

func TestGCPAdapter_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		mockJSONBody   string
		mockStatusCode int
		expectedStatus models.Status
		expectedCount  int
		expectError    bool
	}{
		{
			name:           "Operational status",
			mockJSONBody:    `[]`,
			mockStatusCode: http.StatusOK,
			expectedStatus: models.StatusOperational,
			expectedCount:  0,
		},
		{
			name: "Degraded status",
			mockJSONBody: `[
				{
					"id": "202603240001",
					"service_key": "cloud-sql",
					"service_name": "Cloud SQL",
					"severity": "medium",
					"status_description": "Service Disruption",
					"begin": "2026-03-24T00:01:00Z",
					"end": "",
					"external_desc": "Issue with Cloud SQL"
				}
			]`,
			mockStatusCode: http.StatusOK,
			expectedStatus: models.StatusDegraded,
			expectedCount:  1,
		},
		{
			name:           "HTTP error",
			mockJSONBody:    "Internal Server Error",
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				_, _ = fmt.Fprint(w, tt.mockJSONBody)
			}))
			defer server.Close()

			adapter := &GCPAdapter{BaseURL: server.URL}
			res, err := adapter.Fetch(context.Background())

			if (err != nil) != tt.expectError {
				t.Fatalf("Fetch() error = %v, expectError %v", err, tt.expectError)
			}

			if err == nil {
				if res.Status != tt.expectedStatus {
					t.Errorf("Fetch() status = %v, want %v", res.Status, tt.expectedStatus)
				}
				if len(res.Incidents) != tt.expectedCount {
					t.Errorf("Fetch() incident count = %v, want %v", len(res.Incidents), tt.expectedCount)
				}
				if tt.expectedCount > 0 {
					if res.Incidents[0].ID != "202603240001" {
						t.Errorf("Fetch() incident ID = %v, want 202603240001", res.Incidents[0].ID)
					}
				}
			}
		})
	}
}
