package invoker

import (
	"context"
	"errors"
	"testing"
)

func Test_LogInvoker_InvokedError_ReturnError(t *testing.T) {
	i := NewLogInvoker()
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return errors.New("error")
		},
		MarshalJSONFunc: func() ([]byte, error) {
			return []byte("{\"ExecuteFunc\":null}"), nil
		},
	}
	if err := i.Invoke(context.Background(), c); err == nil {
		t.Error("Expected error")
	}
}

func Test_LogInvoker_InvokedNoError_NoError(t *testing.T) {
	i := NewLogInvoker()
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return nil
		},
		MarshalJSONFunc: func() ([]byte, error) {
			return []byte("{\"ExecuteFunc\":null}"), nil
		},
	}
	if err := i.Invoke(context.Background(), c); err != nil {
		t.Error("Expected nil")
	}
}
