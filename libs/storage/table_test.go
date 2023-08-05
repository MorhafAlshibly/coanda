package storage

import (
	"context"
	"testing"
)

func createTestTable() (*TableStorage, error) {
	store, err := NewTableStorage(context.TODO(), "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", "test")
	if err != nil {
		return store, err
	}
	store.Client.Delete(context.TODO(), nil)
	_, err = store.Client.CreateTable(context.TODO(), nil)
	return store, nil
}

func TestTableStorageAdd(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	key, err := store.Add(context.TODO(), "test", map[string]any{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	if key == "" {
		t.Error("Key not returned")
	}
}

func TestTableStorageGet(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	key, err := store.Add(context.TODO(), "test", map[string]any{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	entity, err := store.Get(context.TODO(), key, "test")
	if err != nil {
		t.Error(err)
	}
	if entity["test"] != "test" {
		t.Error("Wrong data")
	}
}

func TestTableStorageDataDoesNotExist(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	_, err = store.Get(context.TODO(), "test", "test")
	if err == nil {
		t.Error("Error should be thrown")
	}
}

func TestTableStorageQuery(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	key, err := store.Add(context.TODO(), "test", map[string]any{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	entities, err := store.Query(context.TODO(), "RowKey eq '"+key+"'", 1, 1)
	if err != nil {
		t.Error(err)
	}
	if len(entities) != 1 {
		t.Error("Wrong number of entities")
	}
	if entities[0].Data["test"] != "test" {
		t.Error("Wrong data")
	}
}
