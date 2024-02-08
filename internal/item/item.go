package item

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	api.UnimplementedItemServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
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

func (s *Service) GetItem(ctx context.Context, input *api.ItemRequest) (*api.GetItemResponse, error) {
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

func (s *Service) UpdateItem(ctx context.Context, input *api.UpdateItemRequest) (*api.UpdateItemResponse, error) {
	command := NewUpdateItemCommand(s, input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteItem(ctx context.Context, input *api.ItemRequest) (*api.ItemResponse, error) {
	command := NewDeleteItemCommand(s, input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func unmarshalItem(item *model.Item) (*api.Item, error) {
	data, err := conversion.RawJsonToProtobufStruct(item.Data)
	if err != nil {
		return nil, err
	}
	return &api.Item{
		Id:        item.ID,
		Type:      item.Type,
		Data:      data,
		ExpiresAt: timestamppb.New(item.ExpiresAt.Time),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
	}, nil

}

// Enum for errors
type ItemRequestError string

const (
	ID_REQUIRED   ItemRequestError = "ID_REQUIRED"
	TYPE_REQUIRED ItemRequestError = "TYPE_REQUIRED"
)

func (s *Service) checkForItemRequestError(request *api.ItemRequest) *ItemRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_REQUIRED)
	}
	if request.Id == "" {
		return conversion.ValueToPointer(ID_REQUIRED)
	}
	if request.Type == "" {
		return conversion.ValueToPointer(TYPE_REQUIRED)
	}
	return nil
}
