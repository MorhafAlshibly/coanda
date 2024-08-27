package storage

import (
	"context"
	"testing"
)

func Test_MockStorage_Store_FunctionCalled(t *testing.T) {
	m := &MockStorage{
		StoreFunc: func(ctx context.Context, key string, data []byte, metadata map[string]*string) error {
			return nil
		},
	}
	if err := m.Store(context.Background(), "", nil, nil); err != nil {
		t.Error("Expected nil")
	}
}

func Test_MockStorage_Retrieve_FunctionCalled(t *testing.T) {
	m := &MockStorage{
		RetrieveFunc: func(ctx context.Context, key string) ([]byte, error) {
			return nil, nil
		},
	}
	if _, err := m.Retrieve(context.Background(), ""); err != nil {
		t.Error("Expected nil")
	}
}

func Test_MockStorage_Delete_FunctionCalled(t *testing.T) {
	m := &MockStorage{
		DeleteFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}
	if err := m.Delete(context.Background(), ""); err != nil {
		t.Error("Expected nil")
	}
}
