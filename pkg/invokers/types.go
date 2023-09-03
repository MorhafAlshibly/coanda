package invokers

import "context"

type Invoker interface {
	Invoke(ctx context.Context, command Command) error
}

type Command interface {
	Execute(ctx context.Context) error
}
