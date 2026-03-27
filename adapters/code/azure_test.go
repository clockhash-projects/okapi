package code

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"okapi/internal/models"
)

func TestAzureAdapter_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		mockRSSBody    string
		mockStatusCode int
		expectedStatus models.Status
		expectedCount  int
		expectError    bool
	}{
		{
			name: "Operational status (empty feed)",
			mockRSSBody: `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
	<channel>
		<title>Azure Status</title>
	</channel>
</rss>`,
			mockStatusCode: http.StatusOK,
			expectedStatus: models.StatusOperational,
			expectedCount:  0,
		},
		{
			name: "Major outage (active incidents)",
			mockRSSBody: `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
	<channel>
		<item>
			<title>Virtual Machines - East US - Service Issue</title>
			<description>Investigating reports of connectivity issues.</description>
			<pubDate>Mon, 17 Mar 2026 10:00:00 GMT</pubDate>
			<guid>inc123</guid>
		</item>
	</channel>
</rss>`,
			mockStatusCode: http.StatusOK,
			expectedStatus: models.StatusMajorOutage,
			expectedCount:  1,
		},
		{
			name:           "HTTP error",
			mockRSSBody:    "Internal Server Error",
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:           "Invalid XML",
			mockRSSBody:    "not xml",
			mockStatusCode: http.StatusOK,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				_, _ = fmt.Fprint(w, tt.mockRSSBody)
			}))
			defer server.Close()

			adapter := &AzureAdapter{BaseURL: server.URL}
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
					if res.Incidents[0].ID != "inc123" {
						t.Errorf("Fetch() incident ID = %v, want inc123", res.Incidents[0].ID)
					}
				}
			}
		})
	}
}
