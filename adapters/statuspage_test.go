package adapters

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"okapi/internal/models"
)

func TestStatuspageAdapter_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		subdomain      string
		mockStatusURL  string
		mockStatusBody string
		mockStatusCode int
		expected       models.StatusResponse
		expectError    bool
	}{
		{
			name:           "Operational status",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusOK,
			mockStatusBody: `
			{
				"page": { "name": "Test Page", "url": "https://test.statuspage.io" },
				"status": { "indicator": "none", "description": "All systems normal" },
				"components": [
					{ "name": "API", "status": "operational", "updated_at": "2026-03-16T10:00:00Z" }
				],
				"incidents": []
			}
			`,
			expected: models.StatusResponse{
				Service:    "test",
				Status:     models.StatusOperational,
				Summary:    "All systems normal",
				DataSource: "official_api",
				Components: []models.Component{
					{Name: "API", Status: models.StatusOperational, UpdatedAt: time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)},
				},
				Incidents: []models.Incident{},
			},
		},
		{
			name:           "Degraded status",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusOK,
			mockStatusBody: `
			{
				"page": { "name": "Test Page", "url": "https://test.statuspage.io" },
				"status": { "indicator": "minor", "description": "Some services degraded" },
				"components": [
					{ "name": "API", "status": "degraded_performance", "updated_at": "2026-03-16T10:00:00Z" }
				],
				"incidents": []
			}
			`,
			expected: models.StatusResponse{
				Service:    "test",
				Status:     models.StatusDegraded,
				Summary:    "Some services degraded",
				DataSource: "official_api",
				Components: []models.Component{
					{Name: "API", Status: models.StatusDegraded, UpdatedAt: time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)},
				},
				Incidents: []models.Incident{},
			},
		},
		{
			name:           "Major outage status",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusOK,
			mockStatusBody: `
			{
				"page": { "name": "Test Page", "url": "https://test.statuspage.io" },
				"status": { "indicator": "major", "description": "Major outage affecting all systems" },
				"components": [],
				"incidents": [
					{ "id": "inc1", "name": "Core outage", "status": "major", "shortlink": "...", "created_at": "2026-03-16T09:00:00Z", "updated_at": "2026-03-16T09:30:00Z" }
				]
			}
			`,
			expected: models.StatusResponse{
				Service:    "test",
				Status:     models.StatusMajorOutage,
				Summary:    "Major outage affecting all systems",
				DataSource: "official_api",
				Components: []models.Component{},
				Incidents: []models.Incident{
					{ID: "inc1", Title: "Core outage", Status: "major", Body: "...", CreatedAt: time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2026, 3, 16, 9, 30, 0, 0, time.UTC)},
				},
			},
		},
		{
			name:           "Maintenance status",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusOK,
			mockStatusBody: `
			{
				"page": { "name": "Test Page", "url": "https://test.statuspage.io" },
				"status": { "indicator": "under_maintenance", "description": "Scheduled maintenance" },
				"components": [],
				"incidents": []
			}
			`,
			expected: models.StatusResponse{
				Service:    "test",
				Status:     models.StatusMaintenance,
				Summary:    "Scheduled maintenance",
				DataSource: "official_api",
				Components: []models.Component{},
				Incidents:  []models.Incident{},
			},
		},
		{
			name:           "Unknown status",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusOK,
			mockStatusBody: `
			{
				"page": { "name": "Test Page", "url": "https://test.statuspage.io" },
				"status": { "indicator": "unknown", "description": "Status unknown" },
				"components": [],
				"incidents": []
			}
			`,
			expected: models.StatusResponse{
				Service:    "test",
				Status:     models.StatusUnknown,
				Summary:    "Status unknown",
				DataSource: "official_api",
				Components: []models.Component{},
				Incidents:  []models.Incident{},
			},
		},
		{
			name:           "HTTP error",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusInternalServerError,
			mockStatusBody: `{"error": "not found"}`,
			expectError:    true,
		},
		{
			name:           "JSON parse error",
			subdomain:      "test.statuspage.io",
			mockStatusURL:  "/api/v2/summary.json",
			mockStatusCode: http.StatusOK,
			mockStatusBody: `invalid json`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.mockStatusURL {
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				_, err := w.Write([]byte(tt.mockStatusBody))
				if err != nil {
					t.Errorf("failed to write response: %v", err)
				}
			}))
			defer server.Close()

			// Extract subdomain from server URL for the adapter config
			mockSubdomain := strings.TrimPrefix(server.URL, "http://")
			tt.expected.SourceURL = fmt.Sprintf("http://%s/api/v2/summary.json", mockSubdomain)

			cfg := StatuspageConfig{
				ID:                  "test",
				DisplayName:         "Test Service",
				Kind:                "statuspage",
				Subdomain:           mockSubdomain,
				PollIntervalSeconds: 60,
			}
			adapter := NewStatuspageAdapter(cfg)

			res, err := adapter.Fetch(context.Background())

			if (err != nil) != tt.expectError {
				t.Fatalf("Fetch() error = %v, expectError %v", err, tt.expectError)
			}

			if err == nil && !equal(res, &tt.expected) {
				t.Errorf("Fetch() result = %#+v, want %#v", res, &tt.expected)
			}
		})
	}
}

// Helper function to compare StatusResponse structs
func equal(a, b *models.StatusResponse) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.Service != b.Service || a.Status != b.Status || a.Summary != b.Summary || a.DataSource != b.DataSource || a.SourceURL != b.SourceURL || a.Cached != b.Cached {
		return false
	}

	if len(a.Components) != len(b.Components) {
		return false
	}
	for i := range a.Components {
		if a.Components[i].Name != b.Components[i].Name || a.Components[i].Status != b.Components[i].Status || !a.Components[i].UpdatedAt.Equal(b.Components[i].UpdatedAt) {
			return false
		}
	}

	if len(a.Incidents) != len(b.Incidents) {
		return false
	}
	for i := range a.Incidents {
		if a.Incidents[i].ID != b.Incidents[i].ID || a.Incidents[i].Title != b.Incidents[i].Title || a.Incidents[i].Status != b.Incidents[i].Status || a.Incidents[i].Body != b.Incidents[i].Body || !a.Incidents[i].CreatedAt.Equal(b.Incidents[i].CreatedAt) || !a.Incidents[i].UpdatedAt.Equal(b.Incidents[i].UpdatedAt) {
			return false
		}
	}

	return true
}
