package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"

	// // "github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedItemServiceServer
	database             *sqlc.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
	minTypeLength        uint8
	maxTypeLength        uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithDatabase(database *sqlc.Queries) func(*Service) {
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

func WithMinTypeLength(minTypeLength uint8) func(*Service) {
	return func(input *Service) {
		input.minTypeLength = minTypeLength
	}
}

func WithMaxTypeLength(maxTypeLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxTypeLength = maxTypeLength
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
	service := Service{
		minTypeLength:        3,
		maxTypeLength:        20,
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, opt := range opts {
		opt(&service)
	}
	return &service
}

func (s *Service) CreateItem(ctx context.Context, input *api.CreateItemRequest) (*api.CreateItemResponse, error) {
	command := NewCreateItemCommand(s, input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetItem(ctx context.Context, input *api.GetItemRequest) (*api.GetItemResponse, error) {
	command := NewGetItemCommand(s, input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetItems(ctx context.Context, input *api.GetItemsRequest) (*api.GetItemsResponse, error) {
	command := NewGetItemsCommand(s, input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func UnmarshalItem(item *sqlc.Item) (*api.Item, error) {
	data, err := conversion.RawJsonToProtobufStruct(item.Data)
	if err != nil {
		return nil, err
	}
	return &api.Item{
		Id:        item.ID,
		Data:      data,
		ExpiresAt: item.ExpiresAt.Time.Format(time.RFC3339),
	}, nil
}
