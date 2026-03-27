package polling

import (
	"context"
	"testing"
	"time"

	"okapi/adapters"
	"okapi/internal/config"
	"okapi/internal/logger"
	"okapi/internal/models"
	"okapi/internal/webhooks"
)

// Mock items for polling test
type mockAdapter struct {
	id string
}

func (m *mockAdapter) ID() string                  { return m.id }
func (m *mockAdapter) DisplayName() string         { return m.id }
func (m *mockAdapter) PollInterval() time.Duration { return time.Minute }
func (m *mockAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	return &models.StatusResponse{Service: m.id, Status: models.StatusOperational}, nil
}

type mockCache struct {
	getFunc func(key string) (*models.StatusResponse, error)
	setFunc func(key string, val *models.StatusResponse) error
}

func (m *mockCache) Get(ctx context.Context, key string) (*models.StatusResponse, error) {
	if m.getFunc != nil {
		return m.getFunc(key)
	}
	return nil, nil
}
func (m *mockCache) GetAll(ctx context.Context, keys []string) (map[string]*models.StatusResponse, error) {
	return nil, nil
}
func (m *mockCache) Set(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error {
	if m.setFunc != nil {
		return m.setFunc(key, val)
	}
	return nil
}
func (m *mockCache) Close() error { return nil }

func TestWorker_DoPoll(t *testing.T) {
	logger.Init("info", "json")
	ctx := context.Background()
	registry := adapters.NewRegistry()
	adapter := &mockAdapter{id: "test"}
	registry.Register(adapter)

	cacheSetCalled := false
	mCache := &mockCache{
		setFunc: func(key string, val *models.StatusResponse) error {
			cacheSetCalled = true
			return nil
		},
	}

	w := NewWorker(registry, mCache, webhooks.NewManager(), &config.Config{})

	w.doPoll(ctx, adapter)

	if !cacheSetCalled {
		t.Error("expected cache.Set to be called")
	}
}
