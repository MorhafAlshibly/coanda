package storage

import (
	"context"
	"reflect"
	"testing"
)

func createTestTable() (*TableStorage, error) {
	store, err := NewTableStorage(context.TODO(), nil, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", "test")
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
	object, err := store.Add(context.TODO(), "test", map[string]string{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	if object.Key == "" {
		t.Error("Key is empty")
	}
	if object.Pk != "test" {
		t.Error("Wrong partition key")
	}
	if object.Data["test"] != "test" {
		t.Error("Wrong data")
	}
}

func TestTableStorageGet(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	object, err := store.Add(context.TODO(), "test", map[string]string{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	entity, err := store.Get(context.TODO(), object.Key, "test")
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(entity, object) == false {
		t.Error("Wrong data")
	}
}

func TestTableStorageDataDoesNotExist(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	object, err := store.Get(context.TODO(), "test", "test")
	if err == nil {
		t.Error("Error should be thrown")
	}
	if err.Error() != (&ObjectNotFoundError{}).Error() {
		t.Error("Wrong error thrown")
	}
	if object != nil {
		t.Error("Wrong data")
	}
}

func TestTableStorageQuery(t *testing.T) {
	store, err := createTestTable()
	if err != nil {
		t.Error(err)
	}
	object, err := store.Add(context.TODO(), "test", map[string]string{"test": "test"})
	if err != nil {
		t.Error(err)
	}
	entities, err := store.Query(context.TODO(), nil, 1, 1)
	if err != nil {
		t.Error(err)
	}
	if len(entities) != 1 {
		t.Error("Wrong number of entities")
	}
	if reflect.DeepEqual(entities[0], object) == false {
		t.Error("Wrong data")
	}
}
