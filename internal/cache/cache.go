package cache

import (
	"context"
	"time"

	"okapi/internal/models"
)

type Cache interface {
	Get(ctx context.Context, key string) (*models.StatusResponse, error)
	GetAll(ctx context.Context, keys []string) (map[string]*models.StatusResponse, error)
	Set(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error
	Close() error
}
