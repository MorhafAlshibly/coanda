package invokers

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"

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
	log.Info().Str("requestId", requestID).Msg(fmt.Sprintf("Command %T started", command))

	// Start a trace region
	traceFile, err := os.Create(fmt.Sprintf("trace_%s.out", requestID))
	if err != nil {
		log.Error().Err(err).Str("requestId", requestID).Msg("Failed to create trace file")
		return err
	}
	defer traceFile.Close()

	if err := trace.Start(traceFile); err != nil {
		log.Error().Err(err).Str("requestId", requestID).Msg("Failed to start trace")
		return err
	}
	defer trace.Stop()

	// Create a region to trace the invocation
	region := trace.StartRegion(ctx, fmt.Sprintf("Invoke %T", command))
	defer region.End()

	err = i.invoker.Invoke(ctx, command)
	if err != nil {
		log.Error().Err(err).Str("requestId", requestID).Msg(fmt.Sprintf("Command %T errors", command))
		return err
	}

	log.Info().Str("requestId", requestID).Msg(fmt.Sprintf("Command %T executed and output returned", command))
	return nil
}
