package item

import (
	"context"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

func TestGetItem(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return &storage.Object{
				Key:  "test",
				Pk:   "test",
				Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
			}, nil
		},
	}
	service := NewService(WithStore(store))
	c := GetItemCommand{
		service: service,
		In: &api.GetItemRequest{
			Id:   "test",
			Type: "test",
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Item.Id != "test" {
		t.Error("Key not returned")
	}
	if c.Out.Item.Type != "test" {
		t.Error("Wrong type")
	}
	if c.Out.Item.Data["test"] != "test" {
		t.Error("Wrong data")
	}
}

func TestGetItemNotFound(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return nil, &storage.ObjectNotFoundError{}
		},
	}
	service := NewService(WithStore(store))
	c := GetItemCommand{
		service: service,
		In: &api.GetItemRequest{
			Id:   "test",
			Type: "test",
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Item != nil {
		t.Error("Item returned")
	}
	if c.Out.Error != api.GetItemResponse_NOT_FOUND {
		t.Error("Wrong error")
	}
}

func TestGetItemNoId(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return nil, nil
		},
	}
	service := NewService(WithStore(store))
	c := GetItemCommand{
		service: service,
		In: &api.GetItemRequest{
			Id:   "",
			Type: "test",
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Item != nil {
		t.Error("Item returned")
	}
	if c.Out.Error != api.GetItemResponse_ID_NOT_SET {
		t.Error("Wrong error")
	}
}

func TestGetItemNoType(t *testing.T) {
	store := &storage.MockStorage{
		GetFunc: func(ctx context.Context, key string, pk string) (*storage.Object, error) {
			return nil, nil
		},
	}
	service := NewService(WithStore(store))
	c := GetItemCommand{
		service: service,
		In: &api.GetItemRequest{
			Id:   "test",
			Type: "",
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Item != nil {
		t.Error("Item returned")
	}
	if c.Out.Error != api.GetItemResponse_TYPE_TOO_SHORT {
		t.Error("Wrong error")
	}
}
