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
	marshalled, err := sonic.Marshal(item.Data)
	if err != nil {
		return nil, err
	}
	mapData := map[string]string{
		"Type": item.Type,
		"Data": string(marshalled),
	}
	// If the item has an expiry, add it to the map
	if item.Expire != nil {
		mapData["Expire"] = item.Expire.Format(time.RFC3339)
	}
	// Add the item to the store
	object, err := s.store.Add(ctx, item.Type, mapData)
	if err != nil {
		return nil, err
	}
	// Allot the output
	out := &model.Item{
		ID:     object.Key,
		Type:   item.Type,
		Data:   item.Data,
		Expire: item.Expire,
	}
	// Marshal the output to a string to be cached
	marshalled, err = sonic.Marshal(out)
	if err != nil {
		return nil, err
	}
	// Add the item to the cache in a separate thread
	go s.cache.Add(ctx, object.Key, string(marshalled))
	return out, nil
}

// Get is used to get an item
func (s *ItemService) Get(ctx context.Context, item model.GetItem) (*model.Item, error) {
	var out *model.Item
	data, err := s.cache.Get(ctx, item.ID)
	// If the item is not in the cache, get it from the store, else marshal it to output
	if err == nil {
		err = sonic.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		return out, nil
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
	return out, nil
}

// GetAll is used to get all items of a type
func (s *ItemService) GetAll(ctx context.Context, options model.GetItems) ([]*model.Item, error) {
	var items []*storage.Object
	var outs []*model.Item
	// If the type is not nil, set the filter to the type
	filter := ""
	if options.Type != nil {
		filter = "PartitionKey eq '" + *options.Type + "'"
	}
	// Create a key for the cache based on the filter, max and page
	encodedKey := base64.StdEncoding.EncodeToString([]byte(filter + strconv.Itoa(int(*options.Max)) + "|" + strconv.Itoa(*options.Page)))
	data, err := s.cache.Get(ctx, encodedKey)
	// If the data is in the cache, unmarshal it to the output
	if err == nil {
		err = sonic.Unmarshal([]byte(data), &outs)
		if err != nil {
			return nil, err
		}
		return outs, nil
	}
	items, err = s.store.Query(ctx, filter, int32(*options.Max), *options.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		out, err := objectToItem(item)
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
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

func objectToItem(object *storage.Object) (*model.Item, error) {
	var out model.Item
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
