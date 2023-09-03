package invokers

import (
	"context"
)

type BasicInvoker struct {
}

func (i *BasicInvoker) Invoke(ctx context.Context, command Command) error {
	return command.Execute(ctx)
}
