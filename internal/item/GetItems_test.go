package item

import (
	"context"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
)

func TestGetItems(t *testing.T) {
	itemType := "type"
	db := &database.MockDatabase{
		QueryFunc: func(ctx context.Context, input *dynamoTable.QueryInput) ([]map[string]any, error) {
			if input.KeyConditionExpression != "#type = :type" {
				t.Fatal("expected key condition expression to be #type = :type")
			}
			if input.ExpressionAttributeNames["#type"] != "type" {
				t.Fatal("expected expression attribute names to be type")
			}
			if input.ExpressionAttributeValues[":type"] != itemType {
				t.Fatal("expected expression attribute values to be type")
			}
			return []map[string]any{
				{
					"id":       "id",
					"type":     "type",
					"data":     map[string]any{"key": "value"},
					"expireAt": "1970-01-01T00:00:00Z",
				},
			}, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := NewGetItemsCommand(service, &api.GetItemsRequest{
		Type: &itemType,
	})
	err := c.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("expected success to be true")
	}
	if len(c.Out.Items) != 1 {
		t.Fatal("expected items to be 1")
	}
	if c.Out.Items[0].Id != "id" {
		t.Fatal("expected id to be id")
	}
	if c.Out.Items[0].Type != "type" {
		t.Fatal("expected type to be type")
	}
	mapData, err := conversion.ProtobufStructToMap(c.Out.Items[0].Data)
	if err != nil {
		t.Fatal(err)
	}
	if mapData["key"] != "value" {
		t.Fatal("expected data to be data")
	}
	if c.Out.Items[0].ExpireAt.AsTime().Format(time.RFC3339) != "1970-01-01T00:00:00Z" {
		t.Fatal("expected expireAt to be empty")
	}
}

func TestGetItemsNoType(t *testing.T) {
	defaultMaxPageLength := uint8(10)
	db := &database.MockDatabase{
		ScanFunc: func(ctx context.Context, input *dynamoTable.ScanInput) ([]map[string]any, error) {
			if input.ExclusiveStartKey != nil {
				t.Fatal("expected exclusive start key to be nil")
			}
			if input.Max != defaultMaxPageLength {
				t.Fatal("expected max to be 10")
			}
			return []map[string]any{
				{
					"id":       "id",
					"type":     "type",
					"data":     map[string]any{"key": "value"},
					"expireAt": "1970-01-01T00:00:00Z",
				},
			}, nil
		},
	}
	service := NewService(WithDatabase(db), WithDefaultMaxPageLength(defaultMaxPageLength))
	c := NewGetItemsCommand(service, &api.GetItemsRequest{})
	err := c.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("expected success to be true")
	}
	if len(c.Out.Items) != 1 {
		t.Fatal("expected items to be 1")
	}
	if c.Out.Items[0].Id != "id" {
		t.Fatal("expected id to be id")
	}
	if c.Out.Items[0].Type != "type" {
		t.Fatal("expected type to be type")
	}
	mapData, err := conversion.ProtobufStructToMap(c.Out.Items[0].Data)
	if err != nil {
		t.Fatal(err)
	}
	if mapData["key"] != "value" {
		t.Fatal("expected data to be data")
	}
	if c.Out.Items[0].ExpireAt.AsTime().Format(time.RFC3339) != "1970-01-01T00:00:00Z" {
		t.Fatal("expected expireAt to be empty")
	}
}

func TestGetItemsNoItems(t *testing.T) {
	db := &database.MockDatabase{
		ScanFunc: func(ctx context.Context, input *dynamoTable.ScanInput) ([]map[string]any, error) {
			return nil, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := NewGetItemsCommand(service, &api.GetItemsRequest{})
	err := c.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("expected success to be true")
	}
	if len(c.Out.Items) != 0 {
		t.Fatal("expected items to be 0")
	}
	if c.Out.Error != api.GetItemsResponse_NONE {
		t.Fatal("expected error to be NONE")
	}
}
