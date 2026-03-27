package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"okapi/internal/models"
)

var ErrCacheMiss = errors.New("cache miss")

type entry struct {
	val    *models.StatusResponse
	expiry time.Time
}

type MemoryCache struct {
	mu      sync.RWMutex
	entries map[string]entry
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		entries: make(map[string]entry),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (*models.StatusResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.entries[key]
	if !ok || time.Now().After(e.expiry) {
		return nil, ErrCacheMiss
	}

	return e.val, nil
}

func (c *MemoryCache) GetAll(ctx context.Context, keys []string) (map[string]*models.StatusResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*models.StatusResponse)
	now := time.Now()

	for _, k := range keys {
		if e, ok := c.entries[k]; ok && now.Before(e.expiry) {
			result[k] = e.val
		}
	}

	return result, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = entry{
		val:    val,
		expiry: time.Now().Add(ttl),
	}
	return nil
}

func (c *MemoryCache) Close() error {
	return nil
}
