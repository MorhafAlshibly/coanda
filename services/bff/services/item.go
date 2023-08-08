package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/MorhafAlshibly/coanda/libs/cache"
	"github.com/MorhafAlshibly/coanda/libs/storage"
	"github.com/MorhafAlshibly/coanda/services/bff/model"
	"github.com/bytedance/sonic"
)

// ItemService is used to create, get and get all items
type ItemService struct {
	store storage.Storer
	cache cache.Cacher
}

// NewItemService creates a new item service
func NewItemService(store storage.Storer, cache cache.Cacher) *ItemService {
	return &ItemService{
		store: store,
		cache: cache,
	}
}

// Create is used to create a new item
func (s *ItemService) Create(ctx context.Context, item model.CreateItem) (*model.Item, error) {
	var out model.Item
	marshalled, err := sonic.Marshal(item.Data)
	if err != nil {
		return nil, err
	}
	// If the expire is nil, set it to empty time
	if item.Expire == nil {
		defaultTime := time.Unix(0, 0)
		item.Expire = &defaultTime
	}
	// Add the item to the store
	object, err := s.store.Add(ctx, item.Type, map[string]string{
		"Type":   item.Type,
		"Data":   string(marshalled),
		"Expire": fmt.Sprint((*item.Expire).UnixMilli()),
	})
	if err != nil {
		return nil, err
	}
	// Allot the output
	out.ID = object.Key
	out.Type = item.Type
	out.Data = item.Data
	out.Expire = *item.Expire
	// Marshal the output to a string to be cached
	marshalled, err = sonic.Marshal(out)
	if err != nil {
		return nil, err
	}
	// Add the item to the cache in a separate thread
	go s.cache.Add(ctx, object.Key, string(marshalled))
	return &out, nil
}

// Get is used to get an item
func (s *ItemService) Get(ctx context.Context, item model.GetItem) (*model.Item, error) {
	var out model.Item
	data, err := s.cache.Get(ctx, item.ID)
	// If the item is not in the cache, get it from the store, else marshal it to output
	if err == nil {
		err = sonic.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}
	// Get the item from the store
	object, err := s.store.Get(ctx, item.ID, item.Type)
	if err != nil {
		return nil, err
	}
	out, err = objectToItem(object)
	if err != nil {
		return nil, err
	}
	// Marshal the final output to a string to be cached
	marshalled, err := sonic.Marshal(out)
	if err != nil {
		return nil, err
	}
	// Add the item to the cache in a separate thread
	go s.cache.Add(ctx, item.ID, string(marshalled))
	return &out, nil
}

// GetAll is used to get all items of a type
func (s *ItemService) GetAll(ctx context.Context, item model.GetItems) ([]*model.Item, error) {
	var items []storage.Object
	var outs []*model.Item
	var max int32 = 10
	page := 1
	// If the max and page are not nil, set them to the values of the item
	if item.Max != nil {
		max = int32(*item.Max)
	}
	if item.Page != nil {
		page = int(*item.Page)
	}
	// If the type is not nil, set the filter to the type
	var filter string
	if item.Type != nil {
		filter = "PartitionKey eq '" + *item.Type + "'"
	}
	// Create a key for the cache based on the filter, max and page
	encodedKey := base64.StdEncoding.EncodeToString([]byte(filter + "{" + strconv.Itoa(int(max)) + "}" + strconv.Itoa(page)))
	data, err := s.cache.Get(ctx, encodedKey)
	// If the data is in the cache, unmarshal it to the output
	if err == nil {
		err = sonic.Unmarshal([]byte(data), &outs)
		if err != nil {
			return nil, err
		}
		return outs, nil
	}
	items, err = s.store.Query(ctx, filter, max, page)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		out, err := objectToItem(item)
		if err != nil {
			return nil, err
		}
		outs = append(outs, &out)
	}
	// Marshal the final output to a string to be cached
	marshalled, err := sonic.Marshal(outs)
	if err != nil {
		return nil, err
	}
	// Add the item to the cache in a separate thread
	go s.cache.Add(ctx, encodedKey, string(marshalled))
	return outs, nil
}

func objectToItem(object storage.Object) (model.Item, error) {
	var out model.Item
	out.ID = object.Key
	out.Type = object.Data["Type"]
	// Unmarshal to the output
	err := sonic.Unmarshal([]byte(object.Data["Data"]), &out.Data)
	if err != nil {
		return out, err
	}
	millis, err := strconv.ParseInt(object.Data["Expire"], 10, 64)
	if err != nil {
		return out, err
	}
	out.Expire = time.Unix(0, millis*int64(time.Millisecond))
	return out, nil
}
