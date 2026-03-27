package cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
	"okapi/internal/logger"
	"okapi/internal/models"
	"go.uber.org/zap"
)

// TieredCache implements a two-tier caching strategy to optimize for both
// low latency (L1 memory) and shared state (L2 Redis). It uses the singleflight
// pattern to prevent cache stampedes during L2 misses.
type TieredCache struct {
	local    Cache
	remote   Cache
	localTTL time.Duration
	sf       singleflight.Group
}

// NewTieredCache creates a new TieredCache wrapping a local and remote cache.
func NewTieredCache(local, remote Cache, localTTL time.Duration) *TieredCache {
	return &TieredCache{
		local:    local,
		remote:   remote,
		localTTL: localTTL,
	}
}

// Get attempts to retrieve a value from the local cache (L1). If missing,
// it falls back to the remote cache (L2) using singleflight protection.
// Hits on L2 are automatically backfilled into L1.
func (c *TieredCache) Get(ctx context.Context, key string) (*models.StatusResponse, error) {
	// Tier 1: Local Memory (L1) - very fast (~50ns)
	val, err := c.local.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	// Tier 2: Remote Redis (L2) with singleflight protection to prevent stampedes.
	// Only one goroutine per key will execute the provided function; others wait.
	res, err, shared := c.sf.Do(key, func() (interface{}, error) {
		// Double-check L1 inside the singleflight slot to ensure another
		// goroutine didn't already fill it while we waited.
		val, err := c.local.Get(ctx, key)
		if err == nil {
			return val, nil
		}

		val, err = c.remote.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		// Backfill L1 to speed up subsequent requests on this instance.
		if c.localTTL > 0 {
			_ = c.local.Set(ctx, key, val, c.localTTL)
		}
		return val, nil
	})

	if err != nil {
		return nil, err
	}

	if shared {
		logger.Log.Debug("singleflight shared result", zap.String("key", key))
	}

	return res.(*models.StatusResponse), nil
}

// GetAll performs a bulk lookup across both cache tiers.
func (c *TieredCache) GetAll(ctx context.Context, keys []string) (map[string]*models.StatusResponse, error) {
	if len(keys) == 0 {
		return make(map[string]*models.StatusResponse), nil
	}

	// Bulk check L1 first.
	result, _ := c.local.GetAll(ctx, keys)
	
	missingKeys := make([]string, 0)
	for _, k := range keys {
		if _, ok := result[k]; !ok {
			missingKeys = append(missingKeys, k)
		}
	}

	if len(missingKeys) == 0 {
		return result, nil
	}

	// Fetch missing keys from L2.
	remoteResults, err := c.remote.GetAll(ctx, missingKeys)
	if err != nil {
		// If L2 fails, return partial results from L1 instead of a total failure.
		return result, err
	}

	// Merge L2 results and backfill them into L1.
	for k, v := range remoteResults {
		result[k] = v
		if c.localTTL > 0 {
			_ = c.local.Set(ctx, k, v, c.localTTL)
		}
	}

	return result, nil
}

// Set populates both cache tiers. The remote tier uses the provided TTL,
// while the local tier uses the configured localTTL for instance-level freshness.
func (c *TieredCache) Set(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error {
	if c.localTTL > 0 {
		_ = c.local.Set(ctx, key, val, c.localTTL)
	}
	return c.remote.Set(ctx, key, val, ttl)
}

// Close gracefully shuts down both local and remote cache resources.
func (c *TieredCache) Close() error {
	errLocal := c.local.Close()
	errRemote := c.remote.Close()
	
	if errLocal != nil {
		return errLocal
	}
	return errRemote
}
