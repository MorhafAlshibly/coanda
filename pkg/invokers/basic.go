package invokers

import (
	"context"
)

type BasicInvoker struct {
	invoker Invoker
}

func NewBasicInvoker() *BasicInvoker {
	return &BasicInvoker{}
}

func (i *BasicInvoker) SetInvoker(invoker Invoker) *BasicInvoker {
	i.invoker = invoker
	return i
}

func (i *BasicInvoker) Invoke(ctx context.Context, command Command) error {
	return command.Execute(ctx)
}
