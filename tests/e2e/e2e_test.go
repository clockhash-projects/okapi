//go:build e2e
package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// Base URL for the Okapi service
var baseURL = "http://localhost:8081"
var apiKey = "test-secret-key"

func TestMain(m *testing.M) {
	if v := os.Getenv("OKAPI_BASE_URL"); v != "" {
		baseURL = v
	}
	if v := os.Getenv("OKAPI_TEST_API_KEY"); v != "" {
		apiKey = v
	}

	// Wait for service to be ready
	waitForService()

	os.Exit(m.Run())
}

func waitForService() {
	client := &http.Client{Timeout: 2 * time.Second}
	for i := 0; i < 30; i++ {
		resp, err := client.Get(baseURL + "/_health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Service at %s not ready after 30s\n", baseURL)
	os.Exit(1)
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(baseURL + "/_health")
	if err != nil {
		t.Fatalf("Failed to call /_health: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", body["status"])
	}
}

func TestMetrics(t *testing.T) {
	resp, err := http.Get(baseURL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to call /metrics: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestProtectedEndpoint_Unauthorized(t *testing.T) {
	// /services is protected when auth is enabled
	resp, err := http.Get(baseURL + "/services")
	if err != nil {
		t.Fatalf("Failed to call /services: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d. Make sure OKAPI_AUTH_ENABLED=true in the test environment.", resp.StatusCode)
	}
}

func TestProtectedEndpoint_Authorized(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/services", nil)
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call /services: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
