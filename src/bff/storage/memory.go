package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Memory storage

type MemoryStorage struct {
	Container []map[string]any
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Container: []map[string]any{},
	}
}

func (s *MemoryStorage) Add(ctx context.Context, pk string, data map[string]any) (string, error) {
	key := uuid.New().String()
	data["ID"] = key
	s.Container = append(s.Container, data)
	return key, nil
}

func (s *MemoryStorage) Get(ctx context.Context, key string, pk string) (map[string]any, error) {
	for _, item := range s.Container {
		if item["ID"] == key {
			return item, nil
		}
	}
	return nil, errors.New("Data not found")
}

func (s *MemoryStorage) Query(ctx context.Context, filter string, max int32, page int) ([]QueryResult, error) {
	return nil, nil
}

// Memory cache

type MemoryCache struct {
	Container map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		Container: map[string]string{},
	}
}

func (s *MemoryCache) Add(ctx context.Context, key string, data string) error {
	_, ok := s.Container[key]
	if ok {
		return errors.New("Key already exists")
	}
	s.Container[key] = data
	return nil
}

func (s *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	_, ok := s.Container[key]
	if !ok {
		return "", errors.New("Key not found")
	}
	return s.Container[key], nil
}
