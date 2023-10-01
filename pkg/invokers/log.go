package invokers

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogInvoker struct {
	invoker Invoker
}

func NewLogInvoker() *LogInvoker {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &LogInvoker{
		invoker: &BasicInvoker{},
	}
}

func (i *LogInvoker) SetInvoker(invoker Invoker) *LogInvoker {
	i.invoker = invoker
	return i
}

func (i *LogInvoker) Invoke(ctx context.Context, command Command) error {
	requestId := uuid.New().String()
	ctx = context.WithValue(ctx, "requestId", requestId)
	log.Info().Str("requestId", requestId).Msg(fmt.Sprintf("Command %T started", command))
	err := i.invoker.Invoke(ctx, command)
	if err != nil {
		log.Error().Err(err).Str("requestId", requestId).Msg(fmt.Sprintf("Command %T errors", command))
		return err
	}
	log.Info().Str("requestId", requestId).Msg(fmt.Sprintf("Command %T executed and output returned", command))
	return nil
}
