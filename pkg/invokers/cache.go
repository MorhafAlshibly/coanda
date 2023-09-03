package invokers

import (
	"context"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/bytedance/sonic"
)

type CacheInvoker struct {
	cache cache.Cacher
}

func NewCacheInvoker(cache cache.Cacher) *CacheInvoker {
	return &CacheInvoker{cache: cache}
}

func (i *CacheInvoker) Invoke(ctx context.Context, command Command) error {
	key, err := sonic.Marshal(command)
	if err != nil {
		return err
	}
	result, err := i.cache.Get(ctx, string(key))
	if err != nil {
		err = command.Execute(ctx)
		if err != nil {
			return err
		}
		val, err := sonic.Marshal(command)
		if err != nil {
			return err
		}
		return i.cache.Add(ctx, string(key), string(val))
	}
	err = sonic.Unmarshal([]byte(result), command)
	if err != nil {
		return err
	}
	return nil
}
