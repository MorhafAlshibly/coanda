package storage

import "context"

type Storer interface {
	Store(ctx context.Context, key string, data []byte, metadata map[string]*string) error
	Retrieve(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

type MockStorage struct {
	StoreFunc    func(ctx context.Context, key string, data []byte, metadata map[string]*string) error
	RetrieveFunc func(ctx context.Context, key string) ([]byte, error)
	DeleteFunc   func(ctx context.Context, key string) error
}

func (m *MockStorage) Store(ctx context.Context, key string, data []byte, metadata map[string]*string) error {
	return m.StoreFunc(ctx, key, data, metadata)
}

func (m *MockStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	return m.RetrieveFunc(ctx, key)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	return m.DeleteFunc(ctx, key)
}
