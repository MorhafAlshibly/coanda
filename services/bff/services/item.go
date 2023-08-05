package services

import (
	"context"
	"encoding/base64"
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
	// Specify the output
	var out model.Item
	// Marshal the data to a string
	marshalled, err := sonic.Marshal(item.Data)
	if err != nil {
		return nil, err
	}
	// If the expire is nil, set it to empty time
	if item.Expire == nil {
		item.Expire = &time.Time{}
	}
	// Add the item to the store
	key, err := s.store.Add(ctx, item.Type, map[string]any{
		"Type":   item.Type,
		"Data":   string(marshalled),
		"Expire": *item.Expire,
	})
	if err != nil {
		return nil, err
	}
	// Allot the output
	out.ID = key
	out.Type = item.Type
	out.Data = item.Data
	out.Expire = *item.Expire
	// Marshal the output to a string to be cached
	marshalled, err = sonic.Marshal(out)
	if err != nil {
		return nil, err
	}
	// Add the item to the cache in a separate thread
	go s.cache.Add(ctx, key, string(marshalled))
	return &out, nil
}

// Get is used to get an item
func (s *ItemService) Get(ctx context.Context, item model.GetItem) (*model.Item, error) {
	// Specify the output
	var out model.Item
	// Get the item from the cache
	data, err := s.cache.Get(ctx, item.ID)
	// If the item is not in the cache, get it from the store
	if err == nil {
		// If the item is in the cache, unmarshal it to the output
		err = sonic.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		// Set the ID of the output to the ID of the item and return it
		out.ID = item.ID
		return &out, nil
	}
	// Get the item from the store
	dataMap, err := s.store.Get(ctx, item.ID, item.Type)
	if err != nil {
		return nil, err
	}
	// Allot the output
	out.ID = item.ID
	out.Type = item.Type
	out.Expire = dataMap["Expire"].(time.Time)
	// Unmarshal the data to the output
	err = sonic.Unmarshal([]byte(dataMap["Data"].(string)), &out.Data)
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
	// Specify the output
	var items []storage.QueryResult
	var outs []*model.Item
	// Set default values
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
	// If the data is not in the cache, get it from the store
	items, err = s.store.Query(ctx, filter, max, page)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		// Specify the output
		var out model.Item
		// Unmarshal the data to the output
		err = sonic.Unmarshal([]byte(item.Data["Data"].(string)), &out.Data)
		if err != nil {
			return nil, err
		}
		// Allot the output
		out.ID = item.Key
		out.Type = item.Data["Type"].(string)
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
