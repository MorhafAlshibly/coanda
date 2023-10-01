package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	api.UnimplementedItemServiceServer
	store                storage.Storer
	cache                cache.Cacher
	metrics              metrics.Metrics
	defaultMaxPageLength uint64
}

type NewServiceInput struct {
	Store                storage.Storer
	Cache                cache.Cacher
	Metrics              metrics.Metrics
	DefaultMaxPageLength uint64
}

func NewService(input *NewServiceInput) *Service {
	return &Service{
		store:                input.Store,
		cache:                input.Cache,
		metrics:              input.Metrics,
		defaultMaxPageLength: input.DefaultMaxPageLength,
	}
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
	err := sonic.Unmarshal([]byte(object.Data["Data"]), &out.Data)
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
	expire, err := time.Parse(time.RFC3339, object.Data["Expire"])
	if err != nil {
		return nil, err
	}
	out.Expire = timestamppb.New(expire)
	return &out, nil
}
