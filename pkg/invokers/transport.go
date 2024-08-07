package invokers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransportInvoker struct {
	invoker Invoker
}

func NewTransportInvoker() *TransportInvoker {
	return &TransportInvoker{}
}

func (i *TransportInvoker) SetInvoker(invoker Invoker) *TransportInvoker {
	i.invoker = invoker
	return i
}

func (i *TransportInvoker) Invoke(ctx context.Context, command Command) error {
	err := command.Execute(ctx)
	// Small hack to return the error as a gRPC error
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
