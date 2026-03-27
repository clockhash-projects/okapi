package cache

import (
	"context"
	"os"
	"testing"
	"time"

	"okapi/internal/models"
)

func TestRedisCache(t *testing.T) {
	url := os.Getenv("OKAPI_CACHE_REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}

	c, err := NewRedisCache(url)
	if err != nil {
		t.Skipf("skipping redis test: %v", err)
	}
	defer c.Close()

	ctx := context.Background()
	key := "test-redis-key"
	val := &models.StatusResponse{Service: "test-redis"}

	tests := []struct {
		name    string
		action  func() error
		wantErr bool
	}{
		{
			name: "Set value",
			action: func() error {
				return c.Set(ctx, key, val, 1*time.Minute)
			},
		},
		{
			name: "Get value",
			action: func() error {
				got, err := c.Get(ctx, key)
				if err != nil {
					return err
				}
				if got.Service != val.Service {
					t.Errorf("expected %s, got %s", val.Service, got.Service)
				}
				return nil
			},
		},
		{
			name: "Get non-existent",
			action: func() error {
				_, err := c.Get(ctx, "non-existent")
				if err != ErrCacheMiss {
					t.Errorf("expected ErrCacheMiss, got %v", err)
				}
				return nil
			},
		},
		{
			name: "GetAll",
			action: func() error {
				results, err := c.GetAll(ctx, []string{key, "non-existent"})
				if err != nil {
					return err
				}
				if results[key] == nil || results[key].Service != val.Service {
					t.Error("expected valid result for existing key")
				}
				if _, ok := results["non-existent"]; ok {
					t.Error("expected no result for non-existent key")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.action(); (err != nil) != tt.wantErr {
				t.Errorf("%s: action() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
