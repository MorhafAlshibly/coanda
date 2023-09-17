package cache

import "context"

// Cacher is used to cache data
type Cacher interface {
	Add(ctx context.Context, key string, data string) error
	Get(ctx context.Context, key string) (string, error)
}

type MockCache struct {
	AddFunc func(ctx context.Context, key string, data string) error
	GetFunc func(ctx context.Context, key string) (string, error)
}

func (c *MockCache) Add(ctx context.Context, key string, data string) error {
	return c.AddFunc(ctx, key, data)
}

func (c *MockCache) Get(ctx context.Context, key string) (string, error) {
	return c.GetFunc(ctx, key)
}
