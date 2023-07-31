package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/MorhafAlshibly/coanda/src/bff/model"
	"github.com/MorhafAlshibly/coanda/src/bff/storage"
)

type ItemService struct {
	store storage.Storer
	cache storage.Cacher
}

func NewItemService(store storage.Storer, cache storage.Cacher) *ItemService {
	return &ItemService{
		store: store,
		cache: cache,
	}
}

func (s *ItemService) Create(ctx context.Context, item model.CreateItem) (*model.Item, error) {
	var out model.Item
	marshalled, err := json.Marshal(item.Data)
	if err != nil {
		return nil, err
	}
	key, err := s.store.Add(ctx, item.Type, map[string]any{
		"Type": item.Type,
		"Data": string(marshalled),
	})
	if err != nil {
		return nil, err
	}
	out.ID = key
	out.Type = item.Type
	out.Data = item.Data
	marshalled, err = json.Marshal(out)
	err = s.cache.Add(ctx, key, string(marshalled))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *ItemService) Get(ctx context.Context, item model.GetItem) (*model.Item, error) {
	var out model.Item
	data, err := s.cache.Get(ctx, item.ID)
	if err == nil {
		err = json.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		out.ID = item.ID
		return &out, nil
	}
	dataMap, err := s.store.Get(ctx, item.ID, item.Type)
	if err != nil {
		return nil, err
	}
	out.ID = item.ID
	out.Type = item.Type
	err = json.Unmarshal([]byte(dataMap["Data"].(string)), &out.Data)
	if err != nil {
		return nil, err
	}
	marshalled, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	err = s.cache.Add(ctx, item.ID, string(marshalled))
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *ItemService) GetAll(ctx context.Context, item model.GetItems) ([]*model.Item, error) {
	var items []storage.QueryResult
	var outs []*model.Item
	var max int32 = 10
	page := 1
	if item.Max != nil {
		max = int32(*item.Max)
	}
	if item.Page != nil {
		page = int(*item.Page)
	}
	var filter string
	if item.Type != nil {
		filter = "PartitionKey eq '" + *item.Type + "'"
	}
	encodedKey := base64.StdEncoding.EncodeToString([]byte(filter + "{" + strconv.Itoa(int(max)) + "}" + strconv.Itoa(page)))
	data, err := s.cache.Get(ctx, encodedKey)
	if err == nil {
		err = json.Unmarshal([]byte(data), &outs)
		if err != nil {
			return nil, err
		}
		return outs, nil
	} else {
		items, err = s.store.Query(ctx, filter, max, page)
		if err != nil {
			return nil, err
		}
	}
	for _, item := range items {
		var out model.Item
		err = json.Unmarshal([]byte(item.Data["Data"].(string)), &out.Data)
		if err != nil {
			return nil, err
		}
		out.ID = item.Key
		out.Type = item.Data["Type"].(string)
		outs = append(outs, &out)
	}
	marshalled, err := json.Marshal(outs)
	if err != nil {
		return nil, err
	}
	err = s.cache.Add(ctx, encodedKey, string(marshalled))
	if err != nil {
		return nil, err
	}
	return outs, nil
}
