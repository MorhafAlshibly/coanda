package invoker

import (
	"context"
	"testing"
)

func Test_BasicInvoker_Invoked_Executed(t *testing.T) {
	i := NewBasicInvoker()
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return nil
		},
	}
	if err := i.Invoke(context.Background(), c); err != nil {
		t.Error("Expected nil")
	}
}
