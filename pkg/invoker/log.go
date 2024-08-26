package invoker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey string

const (
	requestIDKey contextKey = "requestId"
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
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	input, err := json.Marshal(command)
	if err != nil {
		log.Error().Err(err).Str("requestId", requestID).Msg(fmt.Sprintf("Command %T errors", command))
		return err
	}
	log.Info().Str("requestId", requestID).Str("input", string(input)).Msg(fmt.Sprintf("Command %T started", command))
	err = i.invoker.Invoke(ctx, command)
	if err != nil {
		log.Error().Err(err).Str("requestId", requestID).Msg(fmt.Sprintf("Command %T errors", command))
		return err
	}
	output, err := json.Marshal(command)
	if err != nil {
		log.Info().Str("requestId", requestID).Msg(fmt.Sprintf("Command %T executed and output returned", command))
		return nil
	}
	log.Info().Str("requestId", requestID).Str("output", string(output)).Msg(fmt.Sprintf("Command %T executed and output returned", command))
	return nil
}
