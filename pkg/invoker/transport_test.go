package invoker

import (
	"context"
	"errors"
	"testing"
)

func Test_TransportInvoker_InvokedError_ReturnInternalError(t *testing.T) {
	i := NewTransportInvoker()
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return errors.New("error")
		},
	}
	err := i.Invoke(context.Background(), c)
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "rpc error: code = Internal desc = error" {
		t.Error("Expected rpc error: code = Internal desc = error but got", err)
	}
}

func Test_TransportInvoker_Invoked_ReturnNil(t *testing.T) {
	i := NewTransportInvoker()
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return nil
		},
	}
	err := i.Invoke(context.Background(), c)
	if err != nil {
		t.Error("Expected nil")
	}
}
