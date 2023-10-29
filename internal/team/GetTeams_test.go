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

func TestGetTeams(t *testing.T) {
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
	max := uint32(1)
	page := uint64(1)
	c := GetTeamsCommand{
		service: service,
		In: &api.GetTeamsRequest{
			Max:  &max,
			Page: &page,
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
	if c.Out.Teams == nil {
		t.Error("No record returned")
	}
	if c.Out.Teams[0].Id != id.Hex() {
		t.Error("Wrong id returned")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Error("Wrong name returned")
	}
	if c.Out.Teams[0].Owner != 1 {
		t.Error("Wrong userId returned")
	}
	if c.Out.Teams[0].MembersWithoutOwner[0] != 1 {
		t.Error("Wrong userId returned")
	}
	if c.Out.Teams[0].Score != 0 {
		t.Error("Wrong score returned")
	}
	if c.Out.Teams[0].Data["test"] != "test" {
		t.Error("Wrong data returned")
	}
}
