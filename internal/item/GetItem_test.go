package item

// import (
// 	"context"
// 	"reflect"
// 	"testing"

// 	"github.com/MorhafAlshibly/coanda/api"
// 	"github.com/MorhafAlshibly/coanda/pkg/conversion"
// 	"github.com/MorhafAlshibly/coanda/pkg/database"
// 	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
// )

// func TestGetItem(t *testing.T) {
// 	data := map[string]any{"key": "value"}
// 	db := &database.MockDatabase{
// 		GetItemFunc: func(ctx context.Context, input *dynamoTable.GetItemInput) (map[string]any, error) {
// 			if input.Key["id"] != "id" {
// 				t.Fatal("expected id to be id")
// 			}
// 			if input.Key["type"] != "type" {
// 				t.Fatal("expected type to be type")
// 			}
// 			if input.ProjectionExpression != "" {
// 				t.Fatal("expected projection expression to be empty")
// 			}
// 			return map[string]any{
// 				"id":     "id",
// 				"type":   "type",
// 				"data":   data,
// 				"expire": "",
// 			}, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := NewGetItemCommand(service, &api.GetItemRequest{
// 		Id:   "id",
// 		Type: "type",
// 	})
// 	err := c.Execute(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Fatal("expected success to be true")
// 	}
// 	if c.Out.Item.Id != "id" {
// 		t.Fatal("expected id to be id")
// 	}
// 	if c.Out.Item.Type != "type" {
// 		t.Fatal("expected type to be type")
// 	}
// 	mapData, err := conversion.ProtobufStructToMap(c.Out.Item.Data)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !reflect.DeepEqual(mapData, data) {
// 		t.Fatal("expected data to be data")
// 	}
// 	if c.Out.Item.Expire != "" {
// 		t.Fatal("expected expire to be empty")
// 	}
// 	if c.Out.Error != api.GetItemResponse_NONE {
// 		t.Fatal("expected error to be NONE")
// 	}
// }

// func TestGetItemNotFound(t *testing.T) {
// 	db := &database.MockDatabase{
// 		GetItemFunc: func(ctx context.Context, input *dynamoTable.GetItemInput) (map[string]any, error) {
// 			return nil, &dynamoTable.ItemNotFoundError{}
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := NewGetItemCommand(service, &api.GetItemRequest{
// 		Id:   "id",
// 		Type: "type",
// 	})
// 	err := c.Execute(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Fatal("expected success to be false")
// 	}
// 	if c.Out.Item != nil {
// 		t.Fatal("expected item to be nil")
// 	}
// 	if c.Out.Error != api.GetItemResponse_NOT_FOUND {
// 		t.Fatal("expected error to be NOT_FOUND")
// 	}
// }

// func TestGetItemNoId(t *testing.T) {
// 	service := NewService()
// 	c := NewGetItemCommand(service, &api.GetItemRequest{})
// 	err := c.Execute(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Fatal("expected success to be false")
// 	}
// 	if c.Out.Item != nil {
// 		t.Fatal("expected item to be nil")
// 	}
// 	if c.Out.Error != api.GetItemResponse_ID_NOT_SET {
// 		t.Fatal("expected error to be ID_NOT_SET")
// 	}
// }

// func TestGetItemNoType(t *testing.T) {
// 	service := NewService()
// 	c := NewGetItemCommand(service, &api.GetItemRequest{
// 		Id: "id",
// 	})
// 	err := c.Execute(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Fatal("expected success to be false")
// 	}
// 	if c.Out.Item != nil {
// 		t.Fatal("expected item to be nil")
// 	}
// 	if c.Out.Error != api.GetItemResponse_TYPE_TOO_SHORT {
// 		t.Fatal("expected error to be TYPE_TOO_SHORT")
// 	}
// }

// func TestGetItemTooLongType(t *testing.T) {
// 	service := NewService(WithMaxTypeLength(3))
// 	c := NewGetItemCommand(service, &api.GetItemRequest{
// 		Id:   "id",
// 		Type: "test",
// 	})
// 	err := c.Execute(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Fatal("expected success to be false")
// 	}
// 	if c.Out.Item != nil {
// 		t.Fatal("expected item to be nil")
// 	}
// 	if c.Out.Error != api.GetItemResponse_TYPE_TOO_LONG {
// 		t.Fatal("expected error to be TYPE_TOO_LONG")
// 	}
// }
