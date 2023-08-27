package cache

import "context"

// Cacher is used to cache data
type Cacher interface {
	Add(ctx context.Context, key string, data string) error
	Get(ctx context.Context, key string) (string, error)
}
