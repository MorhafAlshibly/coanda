package item

import (
	"context"
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type Service struct {
	api.UnimplementedItemServiceServer
	store                storage.Storer
	cache                cache.Cacher
	metrics              metrics.Metrics
	minTypeLength        uint8
	maxTypeLength        uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithStore(store storage.Storer) func(*Service) {
	return func(input *Service) {
		input.store = store
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

func objectToItem(object *storage.Object) (*api.Item, error) {
	var out api.Item
	// Unmarshal to the output
	err := json.Unmarshal([]byte(object.Data["Data"]), &out.Data)
	if err != nil {
		return nil, err
	}
	out.Id = object.Key
	out.Type = object.Data["Type"]
	// If the item has an expiry, add it to the output
	_, ok := object.Data["Expire"]
	if !ok {
		return &out, nil
	}
	out.Expire = object.Data["Expire"]
	return &out, nil
}
