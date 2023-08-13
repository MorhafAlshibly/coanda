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
func (s *RedisCache) Add(ctx context.Context, key string, data string) error {
	return s.Client.Set(ctx, key, data, s.ttl).Err()
}

// Get is used to get data from the cache
func (s *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return s.Client.Get(ctx, key).Result()
}
