package invokers

import (
	"context"
	"testing"
)

func Test_MockInvoker_Invoked_FunctionCalled(t *testing.T) {
	i := &MockInvoker{
		InvokeFunc: func(ctx context.Context, c Command) error {
			return nil
		},
	}
	c := &MockCommand{}
	if err := i.Invoke(context.Background(), c); err != nil {
		t.Error("Expected nil")
	}
}

func Test_MockCommand_Executed_FunctionCalled(t *testing.T) {
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return nil
		},
	}
	if err := c.Execute(context.Background()); err != nil {
		t.Error("Expected nil")
	}
}

func Test_MockCommand_MarshalJSON_FunctionCalled(t *testing.T) {
	c := &MockCommand{
		MarshalJSONFunc: func() ([]byte, error) {
			return nil, nil
		},
	}
	if _, err := c.MarshalJSON(); err != nil {
		t.Error("Expected nil")
	}
}

func Test_MockCommand_UnmarshalJSON_FunctionCalled(t *testing.T) {
	c := &MockCommand{
		UnmarshalJSONFunc: func(c *MockCommand, bytes []byte) error {
			return nil
		},
	}
	if err := c.UnmarshalJSON(nil); err != nil {
		t.Error("Expected nil")
	}
}
