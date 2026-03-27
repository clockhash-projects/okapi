package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"okapi/internal/models"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(url string) (*RedisCache, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (*models.StatusResponse, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var res models.StatusResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *RedisCache) GetAll(ctx context.Context, keys []string) (map[string]*models.StatusResponse, error) {
	if len(keys) == 0 {
		return map[string]*models.StatusResponse{}, nil
	}

	vals, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*models.StatusResponse)
	for i, val := range vals {
		if val == nil {
			continue
		}

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			continue
		}

		var res models.StatusResponse
		if err := json.Unmarshal(data, &res); err == nil {
			result[keys[i]] = &res
		}
	}

	return result, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, val *models.StatusResponse, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}


func (c *RedisCache) Close() error {
	return c.client.Close()
}
