package code

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"okapi/internal/models"
)

func TestSlackAdapter_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		mockJSONBody   string
		mockStatusCode int
		expectedStatus models.Status
		expectedCount  int
		expectError    bool
	}{
		{
			name: "Operational status",
			mockJSONBody: `{
				"status": "ok",
				"date_created": "2026-03-22T19:40:19-07:00",
				"active_incidents": []
			}`,
			mockStatusCode: http.StatusOK,
			expectedStatus: models.StatusOperational,
			expectedCount:  0,
		},
		{
			name: "Degraded status",
			mockJSONBody: `{
				"status": "degraded",
				"active_incidents": [
					{
						"id": 123,
						"title": "Connectivity issues",
						"type": "incident",
						"status": "investigating",
						"date_created": "2026-03-22T19:40:19-07:00",
						"date_updated": "2026-03-22T19:40:19-07:00"
					}
				]
			}`,
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

			adapter := &SlackAdapter{BaseURL: server.URL}
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
					if res.Incidents[0].ID != "123" {
						t.Errorf("Fetch() incident ID = %v, want 123", res.Incidents[0].ID)
					}
				}
			}
		})
	}
}
