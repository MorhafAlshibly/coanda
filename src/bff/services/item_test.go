package services

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/src/bff/model"
	"github.com/MorhafAlshibly/coanda/src/bff/storage"
)

type mockStorage struct{}

func (s *mockStorage) Add(ctx context.Context, pk string, data map[string]any) (string, error) {
	return "test", nil
}

func (s *mockStorage) Get(ctx context.Context, key string, pk string) (map[string]any, error) {
	return map[string]any{
		"Type": "test",
		"Data": map[string]any{
			"test": "test",
		},
		"Expire": time.Time{},
	}, nil
}

func (s *mockStorage) Query(ctx context.Context, filter string, max int32, page int) ([]storage.QueryResult, error) {
	return []storage.QueryResult{
		{
			Key: "test",
			Data: map[string]any{
				"Type": "test",
				"Data": map[string]any{
					"test": "test",
				},
				"Expire": time.Time{},
			},
			Pk: "test",
		},
	}, nil
}

type mockCache struct{}

func (c *mockCache) Add(ctx context.Context, key string, data string) error {
	return nil
}

func (c *mockCache) Get(ctx context.Context, key string) (string, error) {
	return "test", nil
}

func TestItemCreate(t *testing.T) {
	store := &mockStorage{}
	cache := &mockCache{}
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
	if item.ID == "" {
		t.Error("Key not returned")
	}
	if item.Type != "test" {
		t.Error("Wrong type")
	}
	if !reflect.DeepEqual(item.Data, data) {
		t.Error("Wrong data")
	}
	if item.Expire != (time.Time{}) {
		t.Error("Wrong expire")
	}
}

func TestItemGet(t *testing.T) {
	store := &mockStorage{}
	cache := &mockCache{}
	service := NewItemService(store, cache)
	item, err := service.Create(context.TODO(), model.CreateItem{
		Type: "test",
		Data: map[string]any{"test": "test"},
	})
	if err != nil {
		t.Error(err)
	}
	response, err := service.Get(context.TODO(), model.GetItem{
		ID:   item.ID,
		Type: "test",
	})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(item, response) {
		t.Error("Wrong data")
	}
}

func TestItemDoesNotExist(t *testing.T) {
	store := storage.NewMemoryStorage()
	cache := storage.NewMemoryCache()
	service := NewItemService(store, cache)
	_, err := service.Get(context.TODO(), model.GetItem{
		ID:   "test",
		Type: "test",
	})
	if err == nil {
		t.Error("Error should be thrown")
	}
	if err.Error() != "Data not found" {
		t.Error("Wrong error")
	}
}
