package invokers

import (
	"context"
)

type MockInvoker struct {
	InvokeFunc func(ctx context.Context, command Command) error
}

func (i *MockInvoker) Invoke(ctx context.Context, command Command) error {
	return i.InvokeFunc(ctx, command)
}

type MockCommand struct {
	ExecuteFunc     func(ctx context.Context) error
	MarshalJSONFunc func() ([]byte, error)
}

func (c *MockCommand) Execute(ctx context.Context) error {
	return c.ExecuteFunc(ctx)
}

func (c *MockCommand) MarshalJSON() ([]byte, error) {
	return c.MarshalJSONFunc()
}
