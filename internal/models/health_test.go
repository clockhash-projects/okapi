package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestStatusResponse_MarshalJSON(t *testing.T) {
	now := time.Now().UTC()
	
	tests := []struct {
		name     string
		response StatusResponse
		expected string
	}{
		{
			name: "Core fields only",
			response: StatusResponse{
				Service:    "test-service",
				Status:     StatusOperational,
				Summary:    "All systems nominal",
				FetchedAt:  now,
				DataSource: "test",
			},
			expected: `{"service":"test-service","status":"operational","summary":"All systems nominal","fetched_at":"` + now.Format(time.RFC3339Nano) + `","data_source":"test"}`,
		},
		{
			name: "With Metadata (flexible fields)",
			response: StatusResponse{
				Service:    "test-service",
				Status:     StatusOperational,
				Summary:    "All systems nominal",
				FetchedAt:  now,
				DataSource: "test",
				Metadata: map[string]any{
					"region": "us-east-1",
					"load":   0.45,
					"details": map[string]string{
						"version": "1.2.3",
					},
				},
			},
			expected: `{"details":{"version":"1.2.3"},"fetched_at":"` + now.Format(time.RFC3339Nano) + `","load":0.45,"region":"us-east-1","service":"test-service","status":"operational","summary":"All systems nominal","data_source":"test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			// Since map order is non-deterministic, we unmarshal both to maps and compare
			var gotMap, wantMap map[string]any
			if err := json.Unmarshal(got, &gotMap); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expected), &wantMap); err != nil {
				t.Fatalf("failed to unmarshal expected: %v", err)
			}

			if len(gotMap) != len(wantMap) {
				t.Errorf("MarshalJSON() length = %v, want %v. Got: %s", len(gotMap), len(wantMap), string(got))
			}

			for k, v := range wantMap {
				gotVal, ok := gotMap[k]
				if !ok {
					t.Errorf("MarshalJSON() missing key %q", k)
					continue
				}
				
				// Handle time comparison separately if needed, but here it's string-matched
				if fmt.Sprintf("%v", gotVal) != fmt.Sprintf("%v", v) {
					t.Errorf("MarshalJSON() key %q: got %v, want %v", k, gotVal, v)
				}
			}
		})
	}
}
