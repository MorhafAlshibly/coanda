package team

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

func TestGetTeamById(t *testing.T) {
	id := primitive.NewObjectID()
	db := &database.MockDatabase{
		AggregateFunc: func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
			return mongo.NewCursorFromDocuments(bson.A{
				bson.D{
					{Key: "_id", Value: id},
					{Key: "name", Value: "test"},
					{Key: "owner", Value: int64(1)},
					{Key: "membersWithoutOwner", Value: bson.A{int64(1)}},
					{Key: "score", Value: int64(0)},
					{Key: "data", Value: map[string]string{"test": "test"}},
				},
			}, nil, nil)
		},
	}
	service := NewService(WithDatabase(db))
	c := GetTeamCommand{
		service: service,
		In: &api.GetTeamRequest{
			Id: id.Hex(),
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
	if c.Out.Error != api.GetTeamResponse_NONE {
		t.Error("Wrong error")
	}
}