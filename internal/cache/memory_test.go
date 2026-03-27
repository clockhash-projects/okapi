package cache

import (
	"context"
	"testing"
	"time"

	"okapi/internal/models"
)

func TestMemoryCache(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()
	key := "test-key"
	val := &models.StatusResponse{Service: "test"}

	// Test Set
	err := c.Set(ctx, key, val, 1*time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	got, err := c.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Service != val.Service {
		t.Errorf("expected %s, got %s", val.Service, got.Service)
	}

	// Test Expiry
	err = c.Set(ctx, "expired", val, -1*time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	_, err = c.Get(ctx, "expired")
	if err != ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}

	// Test Miss
	_, err = c.Get(ctx, "nonexistent")
	if err != ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}
