package event

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedEventServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
	minEventNameLength   uint8
	maxEventNameLength   uint8
	minRoundNameLength   uint8
	maxRoundNameLength   uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithSql(sql *sql.DB) func(*Service) {
	return func(input *Service) {
		input.sql = sql
	}
}

func WithDatabase(database *model.Queries) func(*Service) {
	return func(input *Service) {
		input.database = database
	}
}

func WithCache(cache cache.Cacher) func(*Service) {
	return func(input *Service) {
		input.cache = cache
	}
}

func WithMetrics(metrics metrics.Metrics) func(*Service) {
	return func(input *Service) {
		input.metrics = metrics
	}
}

func WithMinEventNameLength(minEventNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minEventNameLength = minEventNameLength
	}
}

func WithMaxEventNameLength(maxEventNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxEventNameLength = maxEventNameLength
	}
}

func WithMinRoundNameLength(minRoundNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minRoundNameLength = minRoundNameLength
	}
}

func WithMaxRoundNameLength(maxRoundNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxRoundNameLength = maxRoundNameLength
	}
}

func WithDefaultMaxPageLength(defaultMaxPageLength uint8) func(*Service) {
	return func(input *Service) {
		input.defaultMaxPageLength = defaultMaxPageLength
	}
}

func WithMaxMaxPageLength(maxMaxPageLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxMaxPageLength = maxMaxPageLength
	}
}

func NewService(opts ...func(*Service)) *Service {
	service := &Service{
		minEventNameLength:   3,
		maxEventNameLength:   20,
		minRoundNameLength:   3,
		maxRoundNameLength:   20,
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, opt := range opts {
		opt(service)
	}
	return service
}

func (s *Service) CreateEvent(ctx context.Context, in *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	command := NewCreateEventCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) AddEventResult(ctx context.Context, in *api.AddEventResultRequest) (*api.AddEventResultResponse, error) {
	command := NewAddEventResultCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

// Enum for errors
type EventRequestError string

const (
	NAME_TOO_SHORT      EventRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG       EventRequestError = "NAME_TOO_LONG"
	ID_OR_NAME_REQUIRED EventRequestError = "ID_OR_NAME_REQUIRED"
)

func (s *Service) checkForEventRequestError(request *api.EventRequest) *EventRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_NAME_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.Name == nil {
		return conversion.ValueToPointer(ID_OR_NAME_REQUIRED)
	}
	if len(*request.Name) < int(s.minEventNameLength) {
		return conversion.ValueToPointer(NAME_TOO_SHORT)
	}
	if len(*request.Name) > int(s.maxEventNameLength) {
		return conversion.ValueToPointer(NAME_TOO_LONG)
	}
	return nil
}
