package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{Client: client, ttl: ttl}
}

func (s *RedisCache) Add(ctx context.Context, key string, data string) error {
	err := s.Client.Set(ctx, key, data, s.ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisCache) Get(ctx context.Context, key string) (string, error) {
	data, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
