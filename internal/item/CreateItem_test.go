package item

import (
	"context"
	"reflect"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
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
	service := NewService(WithStore(store))
	data := map[string]string{"test": "test"}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "test",
			Data: data,
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
	if !reflect.DeepEqual(c.Out.Item.Data, data) {
		t.Error("Wrong data")
	}
}

func TestItemCreateNoType(t *testing.T) {
	store := &storage.MockStorage{
		AddFunc: func(ctx context.Context, pk string, data map[string]string) (*storage.Object, error) {
			return nil, nil
		},
	}
	service := NewService(WithStore(store))
	data := map[string]string{"test": "test"}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "",
			Data: data,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Success != false {
		t.Error("Success returned")
	}
	if c.Out.Error != api.CreateItemResponse_TYPE_TOO_SHORT {
		t.Error("Wrong error")
	}
}

func TestItemCreateNoData(t *testing.T) {
	store := &storage.MockStorage{
		AddFunc: func(ctx context.Context, pk string, data map[string]string) (*storage.Object, error) {
			return nil, nil
		},
	}
	service := NewService(WithStore(store))
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "test",
			Data: nil,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Success != false {
		t.Error("Success returned")
	}
	if c.Out.Error != api.CreateItemResponse_DATA_NOT_SET {
		t.Error("Wrong error")
	}
}

func TestItemCreateTooLongType(t *testing.T) {
	store := &storage.MockStorage{
		AddFunc: func(ctx context.Context, pk string, data map[string]string) (*storage.Object, error) {
			return nil, nil
		},
	}
	service := NewService(WithStore(store), WithMaxTypeLength(3))
	data := map[string]string{"test": "test"}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "testte",
			Data: data,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Success != false {
		t.Error("Success returned")
	}
	if c.Out.Error != api.CreateItemResponse_TYPE_TOO_LONG {
		t.Error("Wrong error")
	}
}
