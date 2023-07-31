package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMemoryStorageAdd(t *testing.T) {
	storage := NewMemoryStorage()
	key, err := storage.Add(context.TODO(), "", map[string]any{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	if key == "" {
		t.Error("Key not returned")
	}
	if storage.Container[0]["test"] != "test" {
		t.Error("Data not added to store")
	}
}

func TestMemoryStorageGet(t *testing.T) {
	storage := NewMemoryStorage()
	key, err := storage.Add(context.TODO(), "", map[string]any{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	entity, err := storage.Get(context.TODO(), key, "")
	if err != nil {
		t.Error(err)
	}
	if entity["test"] != "test" {
		t.Error("Wrong data")
	}
}

func TestMemoryStorageDataDoesNotExist(t *testing.T) {
	storage := NewMemoryStorage()
	_, err := storage.Get(context.TODO(), "test", "")
	if err == nil {
		t.Error("Error should be thrown")
	}

}

func TestMemoryCacheAdd(t *testing.T) {
	storage := NewMemoryCache()
	key := uuid.New().String()
	err := storage.Add(context.TODO(), key, "test")
	if err != nil {
		t.Error(err)
	}
	if storage.Container[key] != "test" {
		t.Error("Data not added to store")
	}
}

func TestMemoryCacheGet(t *testing.T) {
	storage := NewMemoryCache()
	key := uuid.New().String()
	err := storage.Add(context.TODO(), key, "test")
	if err != nil {
		t.Error(err)
	}
	entity, err := storage.Get(context.TODO(), key)
	if err != nil {
		t.Error(err)
	}
	if entity != "test" {
		t.Error("Wrong data")
	}
}

func TestMemoryCacheDataAlreadyExists(t *testing.T) {
	storage := NewMemoryCache()
	key := uuid.New().String()
	err := storage.Add(context.TODO(), key, "test")
	if err != nil {
		t.Error(err)
	}
	err = storage.Add(context.TODO(), key, "test")
	if err == nil {
		t.Error("Error should be thrown")
	}
}

func TestMemoryCacheDataDoesNotExist(t *testing.T) {
	storage := NewMemoryCache()
	key := uuid.New().String()
	_, err := storage.Get(context.TODO(), key)
	if err == nil {
		t.Error("Error should be thrown")
	}
}
