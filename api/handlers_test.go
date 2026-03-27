package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"okapi/adapters"
	"okapi/internal/config"
	"okapi/internal/models"
	"okapi/internal/webhooks"
)

// MockCache implements cache.Cache
type MockCache struct {
	GetFunc func(ctx context.Context, key string) (*models.StatusResponse, error)
	SetFunc func(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error
}

func (m *MockCache) Get(ctx context.Context, key string) (*models.StatusResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return nil, nil
}

func (m *MockCache) GetAll(ctx context.Context, keys []string) (map[string]*models.StatusResponse, error) {
	return nil, nil
}

func (m *MockCache) Set(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, val, ttl)
	}
	return nil
}

func (m *MockCache) Close() error { return nil }

func TestSelfHealth(t *testing.T) {
	h := &Handlers{}
	req := httptest.NewRequest("GET", "/_health", nil)
	w := httptest.NewRecorder()

	h.SelfHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %s", resp["status"])
	}
}

func TestListServices(t *testing.T) {
	registry := adapters.NewRegistry()
	// No services registered yet

	h := NewHandlers(registry, &MockCache{}, webhooks.NewManager(), &config.Config{})

	req := httptest.NewRequest("GET", "/services", nil)
	w := httptest.NewRecorder()

	h.ListServices(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp []interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp) != 0 {
		t.Errorf("expected 0 services, got %d", len(resp))
	}
}
