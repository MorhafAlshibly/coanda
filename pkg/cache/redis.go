package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is used to cache data in redis
type RedisCache struct {
	Client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new redis cache
func NewRedisCache(addr string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{Client: client, ttl: ttl}
}

// Add is used to add data to the cache
func (c *RedisCache) Add(ctx context.Context, key string, data string) error {
	return c.Client.Set(ctx, key, data, c.ttl).Err()
}

// Get is used to get data from the cache
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}
