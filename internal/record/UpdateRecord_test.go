package record

import (
	"context"
	"reflect"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRecordUpdate(t *testing.T) {
	var checkUpdate primitive.D
	db := &database.MockDatabase{
		UpdateOneFunc: func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, *mongo.WriteException) {
			checkUpdate = update.(primitive.D)
			return &mongo.UpdateResult{
				MatchedCount:  1,
				ModifiedCount: 1,
			}, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := UpdateRecordCommand{
		service: service,
		In: &api.UpdateRecordRequest{
			Request: &api.GetRecordRequest{
				NameUserId: &api.NameUserId{
					Name:   "test",
					UserId: 1,
				},
			},
			Data: map[string]string{"test": "test"},
		},
	}
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Out.Success != true {
		t.Error("Success not set")
	}
	if c.Out.Error != api.UpdateRecordResponse_NONE {
		t.Error("Error set")
	}
	// Check update
	if !reflect.DeepEqual(checkUpdate, primitive.D{
		{Key: "$set", Value: primitive.D{
			{Key: "data", Value: map[string]string{"test": "test"}},
		}},
	}) {
		t.Error("Update not correct")
	}
}
