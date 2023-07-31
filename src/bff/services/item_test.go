package services

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/MorhafAlshibly/coanda/src/bff/model"
	"github.com/MorhafAlshibly/coanda/src/bff/storage"
)

func TestItemCreate(t *testing.T) {
	store := storage.NewMemoryStorage()
	cache := storage.NewMemoryCache()
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
	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Error(err)
	}
	if cache.Container[item.ID] != string(marshalled) {
		t.Error("Data not added to cache")
	}
}

func TestItemGet(t *testing.T) {
	store := storage.NewMemoryStorage()
	cache := storage.NewMemoryCache()
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
