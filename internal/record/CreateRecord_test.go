package record

// import (
// 	"context"
// 	"testing"

// 	"github.com/MorhafAlshibly/coanda/api"
// 	"github.com/MorhafAlshibly/coanda/pkg/database"
// 	"github.com/MorhafAlshibly/coanda/pkg/invokers"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func TestRecordCreate(t *testing.T) {
// 	id := primitive.NewObjectID()
// 	db := &database.MockDatabase{
// 		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
// 			return id, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := CreateRecordCommand{
// 		service: service,
// 		In: &api.CreateRecordRequest{
// 			Name:   "test",
// 			UserId: 1,
// 			Record: 1,
// 			Data:   map[string]string{"test": "test"},
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Error("Success not set")
// 	}
// 	if c.Out.Id != id.Hex() {
// 		t.Error("Id not returned")
// 	}
// 	if c.Out.Error != api.CreateRecordResponse_NONE {
// 		t.Error("Error set")
// 	}
// }

// func TestRecordCreateNoName(t *testing.T) {
// 	db := &database.MockDatabase{
// 		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
// 			return primitive.NilObjectID, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := CreateRecordCommand{
// 		service: service,
// 		In: &api.CreateRecordRequest{
// 			Name:   "",
// 			UserId: 1,
// 			Record: 1,
// 			Data:   map[string]string{"test": "test"},
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Error("Success returned")
// 	}
// 	if c.Out.Error != api.CreateRecordResponse_NAME_TOO_SHORT {
// 		t.Error("Wrong error")
// 	}
// }

// func TestRecordCreateNoUserId(t *testing.T) {
// 	db := &database.MockDatabase{
// 		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
// 			return primitive.NilObjectID, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := CreateRecordCommand{
// 		service: service,
// 		In: &api.CreateRecordRequest{
// 			Name:   "test",
// 			UserId: 0,
// 			Record: 1,
// 			Data:   map[string]string{"test": "test"},
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Error("Success returned")
// 	}
// 	if c.Out.Error != api.CreateRecordResponse_USER_ID_REQUIRED {
// 		t.Error("Wrong error")
// 	}
// }

// func TestRecordCreateNoRecord(t *testing.T) {
// 	db := &database.MockDatabase{
// 		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
// 			return primitive.NilObjectID, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := CreateRecordCommand{
// 		service: service,
// 		In: &api.CreateRecordRequest{
// 			Name:   "test",
// 			UserId: 1,
// 			Record: 0,
// 			Data:   map[string]string{"test": "test"},
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Error("Success returned")
// 	}
// 	if c.Out.Error != api.CreateRecordResponse_RECORD_REQUIRED {
// 		t.Error("Wrong error")
// 	}
// }

// func TestRecordCreateNoData(t *testing.T) {
// 	id := primitive.NewObjectID()
// 	db := &database.MockDatabase{
// 		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
// 			return id, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	c := CreateRecordCommand{
// 		service: service,
// 		In: &api.CreateRecordRequest{
// 			Name:   "test",
// 			UserId: 1,
// 			Record: 1,
// 			Data:   nil,
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Error("Success returned")
// 	}
// 	if c.Out.Id != id.Hex() {
// 		t.Error("Id not returned")
// 	}
// 	if c.Out.Error != api.CreateRecordResponse_NONE {
// 		t.Error("Wrong error")
// 	}
// }

// func TestRecordNameTooLong(t *testing.T) {
// 	db := &database.MockDatabase{
// 		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
// 			return primitive.NilObjectID, nil
// 		},
// 	}
// 	service := NewService(WithDatabase(db), WithMaxRecordNameLength(3))
// 	c := CreateRecordCommand{
// 		service: service,
// 		In: &api.CreateRecordRequest{
// 			Name:   "test",
// 			UserId: 1,
// 			Record: 1,
// 			Data:   map[string]string{"test": "test"},
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != false {
// 		t.Error("Success returned")
// 	}
// 	if c.Out.Error != api.CreateRecordResponse_NAME_TOO_LONG {
// 		t.Error("Wrong error")
// 	}
// }
