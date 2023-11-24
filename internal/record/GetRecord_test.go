package record

import (
	"context"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRecordGetById(t *testing.T) {
	id := primitive.NewObjectID()
	idHex := id.Hex()
	db := &database.MockDatabase{
		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
			return mongo.NewCursorFromDocuments(bson.A{
				bson.D{
					{Key: "_id", Value: id},
					{Key: "name", Value: "test"},
					{Key: "userId", Value: int64(1)},
					{Key: "record", Value: int64(1)},
					{Key: "rank", Value: int32(1)},
					{Key: "data", Value: map[string]string{"test": "test"}},
				},
			}, nil, nil)
		},
	}
	service := NewService(WithDatabase(db))
	c := GetRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			Id: &idHex,
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Success != true {
		t.Error("Success not returned")
	}
	if c.Out.Error != api.GetRecordResponse_NONE {
		t.Error("Wrong error")
	}
	if c.Out.Record == nil {
		t.Error("No record returned")
	}
	if c.Out.Record.Id != id.Hex() {
		t.Error("Wrong id returned")
	}
	if c.Out.Record.Name != "test" {
		t.Error("Wrong name returned")
	}
	if c.Out.Record.UserId != 1 {
		t.Error("Wrong userId returned")
	}
	if c.Out.Record.Record != 1 {
		t.Error("Wrong record returned")
	}
	if c.Out.Record.Rank != 1 {
		t.Error("Wrong rank returned")
	}
	if c.Out.Record.Data["test"] != "test" {
		t.Error("Wrong data returned")
	}
}

func TestRecordGetByNameUserId(t *testing.T) {
	id := primitive.NewObjectID()
	db := &database.MockDatabase{
		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
			return mongo.NewCursorFromDocuments(bson.A{
				bson.D{
					{Key: "_id", Value: id},
					{Key: "name", Value: "test"},
					{Key: "userId", Value: int64(1)},
					{Key: "record", Value: int64(1)},
					{Key: "rank", Value: int32(1)},
					{Key: "data", Value: map[string]string{"test": "test"}},
				},
			}, nil, nil)
		},
	}
	service := NewService(WithDatabase(db))
	c := GetRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			NameUserId: &api.NameUserId{
				Name:   "test",
				UserId: 1,
			},
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Success != true {
		t.Error("Success not returned")
	}
	if c.Out.Error != api.GetRecordResponse_NONE {
		t.Error("Wrong error")
	}
	if c.Out.Record == nil {
		t.Error("No record returned")
	}
	if c.Out.Record.Id != id.Hex() {
		t.Error("Wrong id returned")
	}
	if c.Out.Record.Name != "test" {
		t.Error("Wrong name returned")
	}
	if c.Out.Record.UserId != 1 {
		t.Error("Wrong userId returned")
	}
	if c.Out.Record.Record != 1 {
		t.Error("Wrong record returned")
	}
	if c.Out.Record.Rank != 1 {
		t.Error("Wrong rank returned")
	}
	if c.Out.Record.Data["test"] != "test" {
		t.Error("Wrong data returned")
	}
}

func TestRecordGetNoId(t *testing.T) {
	db := &database.MockDatabase{}
	service := NewService(WithDatabase(db))
	c := GetRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			Id: nil,
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
	if c.Out.Error != api.GetRecordResponse_INVALID {
		t.Error("Wrong error")
	}
}

func TestRecordGetNoNameUserId(t *testing.T) {
	db := &database.MockDatabase{}
	service := NewService(WithDatabase(db))
	c := GetRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			NameUserId: nil,
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
	if c.Out.Error != api.GetRecordResponse_INVALID {
		t.Error("Wrong error")
	}
}
