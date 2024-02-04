package item

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestItemCreate(t *testing.T) {
	db := &database.MockDatabase{
		PutItemFunc: func(ctx context.Context, input *dynamoTable.PutItemInput) error {
			return nil
		},
	}
	service := NewService(WithDatabase(db))
	data := map[string]any{"test": "test"}
	structData, err := conversion.MapToProtobufStruct(data)
	if err != nil {
		t.Error(err)
	}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "test",
			Data: structData,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err = invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Error != api.CreateItemResponse_NONE {
		t.Error("Wrong error")
	}
	if c.Out.Item.Id == "" {
		t.Error("Key not returned")
	}
	if c.Out.Item.Type != "test" {
		t.Error("Wrong type")
	}
	if !reflect.DeepEqual(c.Out.Item.Data, structData) {
		t.Error("Wrong data")
	}
}

func TestItemCreateNoType(t *testing.T) {
	db := &database.MockDatabase{
		PutItemFunc: func(ctx context.Context, input *dynamoTable.PutItemInput) error {
			return nil
		},
	}
	service := NewService(WithDatabase(db))
	data := map[string]any{"test": "test"}
	structData, err := conversion.MapToProtobufStruct(data)
	if err != nil {
		t.Error(err)
	}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Data: structData,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err = invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Error != api.CreateItemResponse_TYPE_TOO_SHORT {
		t.Error("Wrong error")
	}
}

func TestItemCreateNoData(t *testing.T) {
	db := &database.MockDatabase{
		PutItemFunc: func(ctx context.Context, input *dynamoTable.PutItemInput) error {
			return nil
		},
	}
	service := NewService(WithDatabase(db))
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "test",
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Error != api.CreateItemResponse_NONE {
		t.Error("Wrong error")
	}
	if c.Out.Item.Id == "" {
		t.Error("Key not returned")
	}
	if c.Out.Item.Type != "test" {
		t.Error("Wrong type")
	}
	if c.Out.Item.Data != nil {
		t.Error("Wrong data")
	}
}

func TestItemCreateTooLongType(t *testing.T) {
	db := &database.MockDatabase{
		PutItemFunc: func(ctx context.Context, input *dynamoTable.PutItemInput) error {
			return nil
		},
	}
	service := NewService(WithDatabase(db), WithMinTypeLength(3), WithMaxTypeLength(10))
	data := map[string]any{"test": "test"}
	structData, err := conversion.MapToProtobufStruct(data)
	if err != nil {
		t.Error(err)
	}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "testtesttest",
			Data: structData,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err = invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Error != api.CreateItemResponse_TYPE_TOO_LONG {
		t.Error("Wrong error")
	}
}

func TestItemCreateTooShortType(t *testing.T) {
	db := &database.MockDatabase{
		PutItemFunc: func(ctx context.Context, input *dynamoTable.PutItemInput) error {
			return nil
		},
	}
	service := NewService(WithDatabase(db), WithMinTypeLength(3), WithMaxTypeLength(10))
	data := map[string]any{"test": "test"}
	structData, err := conversion.MapToProtobufStruct(data)
	if err != nil {
		t.Error(err)
	}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type: "te",
			Data: structData,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err = invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Error != api.CreateItemResponse_TYPE_TOO_SHORT {
		t.Error("Wrong error")
	}
}

func TestItemCreateValidExpire(t *testing.T) {
	db := &database.MockDatabase{
		PutItemFunc: func(ctx context.Context, input *dynamoTable.PutItemInput) error {
			return nil
		},
	}
	service := NewService(WithDatabase(db))
	data := map[string]any{"test": "test"}
	structData, err := conversion.MapToProtobufStruct(data)
	if err != nil {
		t.Error(err)
	}
	c := CreateItemCommand{
		service: service,
		In: &api.CreateItemRequest{
			Type:     "test",
			Data:     structData,
			ExpireAt: timestamppb.New(time.Now()),
		},
	}
	invoker := invokers.NewBasicInvoker()
	err = invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Error != api.CreateItemResponse_NONE {
		t.Error("Wrong error")
	}
}
