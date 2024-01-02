package record

// import (
// 	"context"
// 	"testing"

// 	"github.com/MorhafAlshibly/coanda/api"
// 	"github.com/MorhafAlshibly/coanda/pkg/database"
// 	"github.com/MorhafAlshibly/coanda/pkg/invokers"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func TestGetRecords(t *testing.T) {
// 	id := primitive.NewObjectID()
// 	db := &database.MockDatabase{
// 		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
// 			return mongo.NewCursorFromDocuments(bson.A{
// 				bson.D{
// 					{Key: "_id", Value: id},
// 					{Key: "name", Value: "test"},
// 					{Key: "userId", Value: int64(1)},
// 					{Key: "record", Value: int64(1)},
// 					{Key: "rank", Value: int32(1)},
// 					{Key: "data", Value: map[string]string{"test": "test"}},
// 				},
// 			}, nil, nil)
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	max := uint32(1)
// 	page := uint64(1)
// 	c := GetRecordsCommand{
// 		service: service,
// 		In: &api.GetRecordsRequest{
// 			Max:  &max,
// 			Page: &page,
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Error("Success not returned")
// 	}
// 	if c.Out.Error != api.GetRecordsResponse_NONE {
// 		t.Error("Wrong error")
// 	}
// 	if c.Out.Records == nil {
// 		t.Error("No record returned")
// 	}
// 	if c.Out.Records[0].Id != id.Hex() {
// 		t.Error("Wrong id returned")
// 	}
// 	if c.Out.Records[0].Name != "test" {
// 		t.Error("Wrong name returned")
// 	}
// 	if c.Out.Records[0].UserId != 1 {
// 		t.Error("Wrong userId returned")
// 	}
// 	if c.Out.Records[0].Record != 1 {
// 		t.Error("Wrong record returned")
// 	}
// 	if c.Out.Records[0].Rank != 1 {
// 		t.Error("Wrong rank returned")
// 	}
// 	if c.Out.Records[0].Data["test"] != "test" {
// 		t.Error("Wrong data returned")
// 	}
// }

// func TestGetRecordsByName(t *testing.T) {
// 	id := primitive.NewObjectID()
// 	db := &database.MockDatabase{
// 		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
// 			return mongo.NewCursorFromDocuments(bson.A{
// 				bson.D{
// 					{Key: "_id", Value: id},
// 					{Key: "name", Value: "test"},
// 					{Key: "userId", Value: int64(1)},
// 					{Key: "record", Value: int64(1)},
// 					{Key: "rank", Value: int32(1)},
// 					{Key: "data", Value: map[string]string{"test": "test"}},
// 				},
// 			}, nil, nil)
// 		},
// 	}
// 	service := NewService(WithDatabase(db))
// 	max := uint32(1)
// 	page := uint64(1)
// 	name := "test"
// 	c := GetRecordsCommand{
// 		service: service,
// 		In: &api.GetRecordsRequest{
// 			Max:  &max,
// 			Page: &page,
// 			Name: &name,
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Error("Success not returned")
// 	}
// 	if c.Out.Error != api.GetRecordsResponse_NONE {
// 		t.Error("Wrong error")
// 	}
// 	if c.Out.Records == nil {
// 		t.Error("No record returned")
// 	}
// 	if c.Out.Records[0].Id != id.Hex() {
// 		t.Error("Wrong id returned")
// 	}
// 	if c.Out.Records[0].Name != "test" {
// 		t.Error("Wrong name returned")
// 	}
// 	if c.Out.Records[0].UserId != 1 {
// 		t.Error("Wrong userId returned")
// 	}
// 	if c.Out.Records[0].Record != 1 {
// 		t.Error("Wrong record returned")
// 	}
// 	if c.Out.Records[0].Rank != 1 {
// 		t.Error("Wrong rank returned")
// 	}
// 	if c.Out.Records[0].Data["test"] != "test" {
// 		t.Error("Wrong data returned")
// 	}
// }

// func TestGetRecordsByNameTooShort(t *testing.T) {
// 	db := &database.MockDatabase{}
// 	service := NewService(WithDatabase(db))
// 	max := uint32(1)
// 	page := uint64(1)
// 	name := "t"
// 	c := GetRecordsCommand{
// 		service: service,
// 		In: &api.GetRecordsRequest{
// 			Max:  &max,
// 			Page: &page,
// 			Name: &name,
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
// 	if c.Out.Error != api.GetRecordsResponse_NAME_TOO_SHORT {
// 		t.Error("Wrong error")
// 	}
// }

// func TestGetRecordsNoFields(t *testing.T) {
// 	id := primitive.NewObjectID()
// 	db := &database.MockDatabase{
// 		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
// 			return mongo.NewCursorFromDocuments(bson.A{
// 				bson.D{
// 					{Key: "_id", Value: id},
// 					{Key: "name", Value: "test"},
// 					{Key: "userId", Value: int64(1)},
// 					{Key: "record", Value: int64(1)},
// 					{Key: "rank", Value: int32(1)},
// 					{Key: "data", Value: map[string]string{"test": "test"}},
// 				},
// 			}, nil, nil)
// 		},
// 	}
// 	service := NewService(WithDatabase(db), WithDefaultMaxPageLength(1), WithMaxMaxPageLength(1))
// 	c := GetRecordsCommand{
// 		service: service,
// 		In:      &api.GetRecordsRequest{},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Error("Success not returned")
// 	}
// 	if c.Out.Error != api.GetRecordsResponse_NONE {
// 		t.Error("Wrong error")
// 	}
// 	if c.Out.Records == nil {
// 		t.Error("No record returned")
// 	}
// 	if c.Out.Records[0].Id != id.Hex() {
// 		t.Error("Wrong id returned")
// 	}
// 	if c.Out.Records[0].Name != "test" {
// 		t.Error("Wrong name returned")
// 	}
// 	if c.Out.Records[0].UserId != 1 {
// 		t.Error("Wrong userId returned")
// 	}
// 	if c.Out.Records[0].Record != 1 {
// 		t.Error("Wrong record returned")
// 	}
// 	if c.Out.Records[0].Rank != 1 {
// 		t.Error("Wrong rank returned")
// 	}
// 	if c.Out.Records[0].Data["test"] != "test" {
// 		t.Error("Wrong data returned")
// 	}
// }

// func TestGetRecordsLargeMax(t *testing.T) {
// 	db := &database.MockDatabase{
// 		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
// 			return mongo.NewCursorFromDocuments(bson.A{
// 				bson.D{
// 					{Key: "_id", Value: primitive.NewObjectID()},
// 					{Key: "name", Value: "test"},
// 					{Key: "userId", Value: int64(1)},
// 					{Key: "record", Value: int64(1)},
// 					{Key: "rank", Value: int32(1)},
// 					{Key: "data", Value: map[string]string{"test": "test"}},
// 				},
// 			}, nil, nil)
// 		},
// 	}
// 	service := NewService(WithDatabase(db), WithDefaultMaxPageLength(1), WithMaxMaxPageLength(1))
// 	max := uint32(2)
// 	page := uint64(1)
// 	c := GetRecordsCommand{
// 		service: service,
// 		In: &api.GetRecordsRequest{
// 			Max:  &max,
// 			Page: &page,
// 		},
// 	}
// 	invoker := invokers.NewBasicInvoker()
// 	err := invoker.Invoke(context.Background(), &c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if c.Out.Success != true {
// 		t.Error("Success not returned")
// 	}
// 	if c.Out.Error != api.GetRecordsResponse_NONE {
// 		t.Error("Wrong error")
// 	}
// 	if c.Out.Records == nil {
// 		t.Error("No record returned")
// 	}
// 	if len(c.Out.Records) != 1 {
// 		t.Error("Wrong number of records returned")
// 	}
// }
