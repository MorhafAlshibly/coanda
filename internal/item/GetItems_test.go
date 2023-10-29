package item

import (
	"context"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

func TestGetItems(t *testing.T) {
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			return []*storage.Object{
				{
					Key:  "test",
					Pk:   "test",
					Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
				},
			}, nil
		},
	}
	service := NewService(WithStore(store))
	itemType := "test"
	c := GetItemsCommand{
		service: service,
		In: &api.GetItemsRequest{
			Type: &itemType,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if len(c.Out.Items) != 1 {
		t.Error("Wrong number of items")
	}
	if c.Out.Items[0].Id != "test" {
		t.Error("Key not returned")
	}
	if c.Out.Items[0].Type != "test" {
		t.Error("Wrong type")
	}
	if c.Out.Items[0].Data["test"] != "test" {
		t.Error("Wrong data")
	}
}

func TestGetItemsNoType(t *testing.T) {
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			return []*storage.Object{
				{
					Key:  "test",
					Pk:   "test",
					Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
				},
				{
					Key:  "test2",
					Pk:   "test2",
					Data: map[string]string{"Type": "test2", "Data": "{\"test2\":\"test2\"}"},
				},
			}, nil
		},
	}
	service := NewService(WithStore(store))
	c := GetItemsCommand{
		service: service,
		In:      &api.GetItemsRequest{},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if len(c.Out.Items) != 2 {
		t.Error("Wrong number of items")
	}
	if c.Out.Items[0].Id != "test" {
		t.Error("Key not returned")
	}
	if c.Out.Items[0].Type != "test" {
		t.Error("Wrong type")
	}
	if c.Out.Items[0].Data["test"] != "test" {
		t.Error("Wrong data")
	}
	if c.Out.Items[1].Id != "test2" {
		t.Error("Key not returned")
	}
	if c.Out.Items[1].Type != "test2" {
		t.Error("Wrong type")
	}
	if c.Out.Items[1].Data["test2"] != "test2" {
		t.Error("Wrong data")
	}
}

func TestGetItemsNoItems(t *testing.T) {
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			return []*storage.Object{}, nil
		},
	}
	service := NewService(WithStore(store))
	c := GetItemsCommand{
		service: service,
		In:      &api.GetItemsRequest{},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if len(c.Out.Items) != 0 {
		t.Error("Wrong number of items")
	}
}

func TestGetItemsCustomMaxAndPage(t *testing.T) {
	var checkMax int32
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			checkMax = max
			return []*storage.Object{
				{
					Key:  "test",
					Pk:   "test",
					Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
				},
			}, nil
		},
	}
	service := NewService(WithStore(store))
	max := uint32(1)
	page := uint64(1)
	c := GetItemsCommand{
		service: service,
		In: &api.GetItemsRequest{
			Max:  &max,
			Page: &page,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if len(c.Out.Items) != 1 {
		t.Error("Wrong number of items")
	}
	if c.Out.Items[0].Id != "test" {
		t.Error("Key not returned")
	}
	if c.Out.Items[0].Type != "test" {
		t.Error("Wrong type")
	}
	if c.Out.Items[0].Data["test"] != "test" {
		t.Error("Wrong data")
	}
	if checkMax != 1 {
		t.Error("Wrong max")
	}
}

func TestGetItemsLargeMax(t *testing.T) {
	var checkMax int32
	store := &storage.MockStorage{
		QueryFunc: func(ctx context.Context, filter string, max int32, page int) ([]*storage.Object, error) {
			checkMax = max
			return []*storage.Object{
				{
					Key:  "test",
					Pk:   "test",
					Data: map[string]string{"Type": "test", "Data": "{\"test\":\"test\"}"},
				},
			}, nil
		},
	}
	service := NewService(WithStore(store), WithDefaultMaxPageLength(1), WithMaxMaxPageLength(1))
	max := uint32(2)
	page := uint64(1)
	c := GetItemsCommand{
		service: service,
		In: &api.GetItemsRequest{
			Max:  &max,
			Page: &page,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if len(c.Out.Items) != 1 {
		t.Error("Wrong number of items")
	}
	if c.Out.Items[0].Id != "test" {
		t.Error("Key not returned")
	}
	if c.Out.Items[0].Type != "test" {
		t.Error("Wrong type")
	}
	if c.Out.Items[0].Data["test"] != "test" {
		t.Error("Wrong data")
	}
	if checkMax != 1 {
		t.Error("Wrong max")
	}
}
