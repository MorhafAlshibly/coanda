package storage

import (
	"context"
)

type QueryResult struct {
	Key  string
	Data map[string]any
	Pk   string
}

type Storer interface {
	Add(ctx context.Context, pk string, data map[string]any) (string, error)
	Get(ctx context.Context, key string, pk string) (map[string]any, error)
	Query(ctx context.Context, filter string, max int32, page int) ([]QueryResult, error)
}

type Cacher interface {
	Add(ctx context.Context, key string, data string) error
	Get(ctx context.Context, key string) (string, error)
}
