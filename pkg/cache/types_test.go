package cache

import (
	"context"
	"testing"
)

func Test_MockCacheAdd_FunctionReturnNil_NilReturned(t *testing.T) {
	c := &MockCache{
		AddFunc: func(ctx context.Context, key string, data string) error {
			return nil
		},
	}
	if err := c.Add(context.Background(), "key", "data"); err != nil {
		t.Error("Expected nil")
	}
}

func Test_MockCacheGet_FunctionReturnData_DataReturned(t *testing.T) {
	c := &MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "data", nil
		},
	}
	if data, err := c.Get(context.Background(), "key"); err != nil || data != "data" {
		t.Error("Expected data to be 'data'")
	}
}
