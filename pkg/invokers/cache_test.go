package invokers

import (
	"context"
	"errors"
	"testing"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
)

func TestCacheInvokerNotInCache(t *testing.T) {
	cacheGet := false
	invoked := false
	marshalled := false
	c := &cache.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			cacheGet = true
			return "", errors.New("not found")
		},
		AddFunc: func(ctx context.Context, key string, value string) error {
			return nil
		},
	}
	command := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			invoked = true
			return nil
		},
		MarshalJSONFunc: func() ([]byte, error) {
			marshalled = true
			return []byte("{\"ExecuteFunc\":null}"), nil
		},
	}
	invoker := NewCacheInvoker(c)
	err := invoker.Invoke(context.Background(), command)
	if err != nil {
		t.Error(err)
	}
	if !cacheGet {
		t.Error("Expected cache.Get to be called")
	}
	if !invoked {
		t.Error("Expected command.Execute to be called")
	}
	if !marshalled {
		t.Error("Expected command.MarshalJSON to be called")
	}
}

func TestCacheInvokerInCache(t *testing.T) {
	cacheGet := false
	invoked := false
	marshalled := false
	unmarshalled := false
	c := &cache.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			cacheGet = true
			return "{\"ExecuteFunc\":null}", nil
		},
		AddFunc: func(ctx context.Context, key string, value string) error {
			return nil
		},
	}
	command := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			invoked = true
			return nil
		},
		MarshalJSONFunc: func() ([]byte, error) {
			marshalled = true
			return []byte("{\"ExecuteFunc\":null}"), nil
		},
		UnmarshalJSONFunc: func(mockCommand *MockCommand, bytes []byte) error {
			mockCommand.ExecuteFunc = func(ctx context.Context) error {
				return nil
			}
			unmarshalled = true
			return nil
		},
	}
	invoker := NewCacheInvoker(c)
	err := invoker.Invoke(context.Background(), command)
	if err != nil {
		t.Error(err)
	}
	if !cacheGet {
		t.Error("Expected cache.Get to be called")
	}
	if invoked {
		t.Error("Expected command.Execute not to be called")
	}
	if !marshalled {
		t.Error("Expected command.MarshalJSON to be called")
	}
	if !unmarshalled {
		t.Error("Expected command.UnmarshalJSON to be called")
	}
}

func TestCacheInvokerGenerateKey(t *testing.T) {
	key, err := generateKey(&MockCommand{
		ExecuteFunc: nil,
		MarshalJSONFunc: func() ([]byte, error) {
			return []byte("{\"ExecuteFunc\":null}"), nil
		},
	})
	expected := "*invokers.MockCommand: {\"ExecuteFunc\":null}"
	if err != nil {
		t.Error(err)
	}
	if key != expected {
		t.Errorf("Expected key to be %s', got '%s'", expected, key)
	}
}
