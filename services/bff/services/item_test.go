package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/libs/storage"
	"github.com/MorhafAlshibly/coanda/services/bff/model"
)

type mockStorage struct {
	AddResponse   string
	GetResponse   map[string]any
	QueryResponse []storage.QueryResult
	Error         error
}

func (s *mockStorage) Add(ctx context.Context, pk string, data map[string]any) (string, error) {
	return s.AddResponse, s.Error
}

func (s *mockStorage) Get(ctx context.Context, key string, pk string) (map[string]any, error) {
	return s.GetResponse, s.Error
}

func (s *mockStorage) Query(ctx context.Context, filter string, max int32, page int) ([]storage.QueryResult, error) {
	return s.QueryResponse, s.Error
}

type mockCache struct {
	GetResponse string
	Error       error
}

func (c *mockCache) Add(ctx context.Context, key string, data string) error {
	return c.Error
}

func (c *mockCache) Get(ctx context.Context, key string) (string, error) {
	return c.GetResponse, c.Error
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
	store.AddResponse = "test"
	cache.Error = nil
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
	if item.Expire != (time.Time{}) {
		t.Error("Wrong expire")
	}
}

func TestItemGet(t *testing.T) {
	store := &mockStorage{}
	cache := &mockCache{}
	service := NewItemService(store, cache)
	store.GetResponse = map[string]any{
		"Type":   "test",
		"Data":   `{"test":"test"}`,
		"Expire": time.Time{},
	}
	cache.Error = errors.New("Data not found")
	response, err := service.Get(context.TODO(), model.GetItem{
		ID:   "test",
		Type: "test",
	})
	if err != nil {
		t.Error(err)
	}
	if response.ID != "test" {
		t.Error("Wrong key")
	}
	if response.Type != "test" {
		t.Error("Wrong type")
	}
	if response.Expire != (time.Time{}) {
		t.Error("Wrong expire")
	}
	t.Log(response)
}

func TestItemDoesNotExist(t *testing.T) {
	store := &mockStorage{}
	cache := &mockCache{}
	service := NewItemService(store, cache)
	store.GetResponse = nil
	store.Error = errors.New("Data not found")
	cache.Error = errors.New("Data not found")
	_, err := service.Get(context.TODO(), model.GetItem{
		ID:   "nottest",
		Type: "test",
	})
	if err == nil {
		t.Error("Error should be thrown")
	}
	if err.Error() != "Data not found" {
		t.Error("Wrong error")
	}
}
