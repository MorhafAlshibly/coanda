package record

import (
	"context"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestDeleteRecord(t *testing.T) {
	db := &database.MockDatabase{
		DeleteOneFunc: func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException) {
			return &mongo.DeleteResult{
				DeletedCount: 1,
			}, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := DeleteRecordCommand{
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
	if c.Out.Error != api.DeleteRecordResponse_NONE {
		t.Error("Wrong error")
	}
}

func TestDeleteRecordById(t *testing.T) {
	db := &database.MockDatabase{
		DeleteOneFunc: func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException) {
			return &mongo.DeleteResult{
				DeletedCount: 1,
			}, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := DeleteRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			Id: primitive.NewObjectID().Hex(),
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
	if c.Out.Error != api.DeleteRecordResponse_NONE {
		t.Error("Wrong error")
	}
	t.Log(err)
	t.Log(c.Out)
}

func TestDeleteRecordNoIdNoNameUserId(t *testing.T) {
	db := &database.MockDatabase{
		DeleteOneFunc: func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException) {
			return &mongo.DeleteResult{
				DeletedCount: 1,
			}, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := DeleteRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			Id:         "",
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
	if c.Out.Error != api.DeleteRecordResponse_INVALID {
		t.Error("Wrong error")
	}
}

func TestDeleteRecordNoUserId(t *testing.T) {
	db := &database.MockDatabase{
		DeleteOneFunc: func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException) {
			return &mongo.DeleteResult{
				DeletedCount: 1,
			}, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := DeleteRecordCommand{
		service: service,
		In: &api.GetRecordRequest{
			Id: "",
			NameUserId: &api.NameUserId{
				Name:   "test",
				UserId: 0,
			},
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
	if c.Out.Error != api.DeleteRecordResponse_INVALID {
		t.Error("Wrong error")
	}
}
