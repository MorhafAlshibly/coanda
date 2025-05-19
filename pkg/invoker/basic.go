package invoker

import (
	"context"
)

type BasicInvoker struct{}

func NewBasicInvoker() *BasicInvoker {
	return &BasicInvoker{}
}

func (i *BasicInvoker) Invoke(ctx context.Context, command Command) error {
	return command.Execute(ctx)
}
