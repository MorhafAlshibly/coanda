package storage

import (
	"context"
)

// Storer is used to store data
type Storer interface {
	Add(ctx context.Context, pk string, data map[string]string) (*Object, error)
	Get(ctx context.Context, key string, pk string) (*Object, error)
	Query(ctx context.Context, filter string, max int32, page int) ([]*Object, error)
}

// Object is the data structure used to store data
type Object struct {
	Key  string
	Pk   string
	Data map[string]string
}

// Storage errors
type ObjectNotFoundError struct{}

func (e *ObjectNotFoundError) Error() string {
	return "Object not found"
}

type PageNotFoundError struct{}

func (e *PageNotFoundError) Error() string {
	return "Page not found"
}
