package storage

import (
	"context"
)

// QueryResult is used to store the result of a query
type QueryResult struct {
	Key  string
	Pk   string
	Data map[string]any
}

// Storer is used to store data
type Storer interface {
	Add(ctx context.Context, pk string, data map[string]any) (string, error)
	Get(ctx context.Context, key string, pk string) (map[string]any, error)
	Query(ctx context.Context, filter string, max int32, page int) ([]QueryResult, error)
}
