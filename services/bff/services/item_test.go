package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/libs/cache"
	"github.com/MorhafAlshibly/coanda/libs/storage"
	"github.com/MorhafAlshibly/coanda/services/bff/model"
)

func TestItemCreate(t *testing.T) {
	store := &storage.MockStorage{
		AddFunc: func(ctx context.Context, pk string, data map[string]string) (storage.Object, error) {
			return storage.Object{
				Key:  "test",
				Pk:   "test",
				Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}", "Expire": "0"},
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
	if !reflect.DeepEqual(item.Expire, time.Unix(0, 0)) {
		t.Error("Wrong expire")
	}
}

func TestItemGet(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (storage.Object, error) {
			return storage.Object{
				Key:  "test",
				Pk:   "test",
				Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}", "Expire": "0"},
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
	if !reflect.DeepEqual(item.Expire, time.Unix(0, 0)) {
		t.Error("Wrong expire")
	}
}

func TestItemDoesNotExist(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (storage.Object, error) {
			return storage.Object{}, &storage.ObjectNotFoundError{}
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
	t.Log(err)
	t.Log((&storage.ObjectNotFoundError{}).Error())
	if err.Error() != (&storage.ObjectNotFoundError{}).Error() {
		t.Error("Wrong error")
	}
}
