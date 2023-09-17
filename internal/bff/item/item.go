package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api/gql"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/bytedance/sonic"
)

type Service struct {
	store   storage.Storer
	cache   cache.Cacher
	metrics metrics.Metrics
}

func NewService(store storage.Storer, cache cache.Cacher, metrics metrics.Metrics) *Service {
	return &Service{
		store:   store,
		cache:   cache,
		metrics: metrics,
	}
}

func (s *Service) CreateItem(ctx context.Context, input gql.CreateItem) (*gql.Item, error) {
	command := NewCreateItemCommand(s, &input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetItem(ctx context.Context, input gql.GetItem) (*gql.Item, error) {
	command := NewGetItemCommand(s, &input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetItems(ctx context.Context, input gql.GetItems) ([]*gql.Item, error) {
	command := NewGetItemsCommand(s, &input)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func objectToItem(object *storage.Object) (*gql.Item, error) {
	var out gql.Item
	// Unmarshal to the output
	err := sonic.Unmarshal([]byte(object.Data["Data"]), &out.Data)
	if err != nil {
		return nil, err
	}
	out.ID = object.Key
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
	out.Expire = &expire
	return &out, nil
}
