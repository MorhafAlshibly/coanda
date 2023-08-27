package storage

import "context"

// MockStorage is used to mock the storage
type MockStorage struct {
	AddFunc   func(ctx context.Context, pk string, data map[string]string) (*Object, error)
	GetFunc   func(ctx context.Context, key string, pk string) (*Object, error)
	QueryFunc func(ctx context.Context, filter string, max int32, page int) ([]*Object, error)
}

func (s *MockStorage) Add(ctx context.Context, pk string, data map[string]string) (*Object, error) {
	return s.AddFunc(ctx, pk, data)
}

func (s *MockStorage) Get(ctx context.Context, key string, pk string) (*Object, error) {
	return s.GetFunc(ctx, key, pk)
}

func (s *MockStorage) Query(ctx context.Context, filter string, max int32, page int) ([]*Object, error) {
	return s.QueryFunc(ctx, filter, max, page)
}
