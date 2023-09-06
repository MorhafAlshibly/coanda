package invokers

import (
	"context"
	"fmt"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/bytedance/sonic"
)

type CacheInvoker struct {
	invoker Invoker
	cache   cache.Cacher
}

func NewCacheInvoker(cache cache.Cacher) *CacheInvoker {
	return &CacheInvoker{
		invoker: &BasicInvoker{},
		cache:   cache,
	}
}

func (i *CacheInvoker) SetInvoker(invoker Invoker) *CacheInvoker {
	i.invoker = invoker
	return i
}

func (i *CacheInvoker) Invoke(ctx context.Context, command Command) error {
	key, err := generateKey(command)
	if err != nil {
		return err
	}
	result, err := i.cache.Get(ctx, key)
	if err != nil {
		err = i.invoker.Invoke(ctx, command)
		if err != nil {
			return err
		}
		val, err := sonic.Marshal(command)
		if err != nil {
			return err
		}
		go i.cache.Add(context.Background(), key, string(val))
		return nil
	}
	err = sonic.Unmarshal([]byte(result), command)
	if err != nil {
		return err
	}
	return nil
}

func generateKey(command Command) (string, error) {
	commandType := fmt.Sprintf("%T", command)
	commandValue, err := sonic.Marshal(command)
	if err != nil {
		return "", err
	}
	return commandType + ": " + string(commandValue), nil
}
