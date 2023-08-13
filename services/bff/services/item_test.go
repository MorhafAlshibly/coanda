package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/MorhafAlshibly/coanda/services/bff/model"
)

func TestItemCreate(t *testing.T) {
	store := &storage.MockStorage{
		AddFunc: func(ctx context.Context, pk string, data map[string]string) (*storage.Object, error) {
			return &storage.Object{
				Key:  "test",
				Pk:   "test",
				Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
			}, nil
		},
	}
	cache := &cache.MockCache{
		AddFunc: func(ctx context.Context, key string, data string) error {
			return nil
		},
	}
	service := NewItemService(store, cache)
	data := map[string]any{"test": "test"}
	createItem := model.CreateItem{
		Type: "test",
		Data: data,
	}
	item, err := service.Create(context.TODO(), createItem)
	if err != nil {
		t.Error(err)
	}
	if item.ID != "test" {
		t.Error("Key not returned")
	}
	if item.Type != "test" {
		t.Error("Wrong type")
	}
	if !reflect.DeepEqual(item.Data, data) {
		t.Error("Wrong data")
	}
}

func TestItemGet(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return &storage.Object{
				Key:  "test",
				Pk:   "test",
				Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
			}, nil
		},
	}
	cache := &cache.MockCache{
		AddFunc: func(ctx context.Context, key string, data string) error {
			return nil
		},
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("404")
		},
	}
	service := NewItemService(store, cache)
	item, err := service.Get(context.TODO(), model.GetItem{
		ID:   "test",
		Type: "test",
	})
	if err != nil {
		t.Error(err)
	}
	if item.ID != "test" {
		t.Error("Wrong key")
	}
	if item.Type != "test" {
		t.Error("Wrong type")
	}
	if !reflect.DeepEqual(item.Data, map[string]any{"test": "test"}) {
		t.Error("Wrong data")
	}
}

func TestItemDoesNotExist(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return nil, &storage.ObjectNotFoundError{}
		},
	}
	cache := &cache.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("404")
		},
	}
	service := NewItemService(store, cache)
	_, err := service.Get(context.TODO(), model.GetItem{
		ID:   "test",
		Type: "test",
	})
	if err == nil {
		t.Error("Error should be thrown")
	}
	if err.Error() != (&storage.ObjectNotFoundError{}).Error() {
		t.Error("Wrong error")
	}
}

func TestItemGetFromCache(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return nil, &storage.ObjectNotFoundError{}
		},
	}
	cache := &cache.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "{\"ID\":\"test\",\"Type\":\"test\",\"Data\":{\"test\":\"test\"}}", nil
		},
	}
	service := NewItemService(store, cache)
	item, err := service.Get(context.TODO(), model.GetItem{
		ID:   "test",
		Type: "test",
	})
	if err != nil {
		t.Error(err)
	}
	if item.ID != "test" {
		t.Error("Wrong key")
	}
	if item.Type != "test" {
		t.Error("Wrong type")
	}
	if !reflect.DeepEqual(item.Data, map[string]any{"test": "test"}) {
		t.Error("Wrong data")
	}
}

func TestItemGetAll(t *testing.T) {
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			var out []*storage.Object
			for i := 0; i < 10; i++ {
				out = append(out, &storage.Object{
					Key:  "test" + fmt.Sprint(i),
					Pk:   "test",
					Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
				})
			}
			return out, nil
		},
	}
	cache := &cache.MockCache{
		AddFunc: func(ctx context.Context, key string, data string) error {
			return nil
		},
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("404")
		},
	}
	service := NewItemService(store, cache)
	max := int(10)
	page := int(1)
	items, err := service.GetAll(context.TODO(), model.GetItems{
		Max:  &max,
		Page: &page,
	})
	if err != nil {
		t.Error(err)
	}
	if len(items) != 10 {
		t.Error("Wrong length")
	}
	for i := 0; i < 10; i++ {
		if items[i].ID != "test"+fmt.Sprint(i) {
			t.Error("Wrong key")
		}
		if items[i].Type != "test" {
			t.Error("Wrong type")
		}
		if !reflect.DeepEqual(items[i].Data, map[string]any{"test": "test"}) {
			t.Error("Wrong data")
		}
	}
}

func TestItemGetAllFromCache(t *testing.T) {
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			return nil, nil
		},
	}
	cache := &cache.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "[{\"ID\":\"test\",\"Type\":\"test\",\"Data\":{\"test\":\"test\"}}]", nil
		},
	}
	service := NewItemService(store, cache)
	max := int(10)
	page := int(1)
	items, err := service.GetAll(context.TODO(), model.GetItems{
		Max:  &max,
		Page: &page,
	})
	if err != nil {
		t.Error(err)
	}
	if len(items) != 1 {
		t.Error("Wrong length")
	}
	if items[0].ID != "test" {
		t.Error("Wrong key")
	}
	if items[0].Type != "test" {
		t.Error("Wrong type")
	}
	if !reflect.DeepEqual(items[0].Data, map[string]any{"test": "test"}) {
		t.Error("Wrong data")
	}
}
