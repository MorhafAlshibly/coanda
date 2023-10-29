package team

import (
	"context"
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateTeam(t *testing.T) {
	id := primitive.NewObjectID()
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return id, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "test",
			Owner:               1,
			MembersWithoutOwner: []uint64{},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NONE {
		t.Error("Wrong error")
	}
	if c.Out.Id != id.Hex() {
		t.Error("Wrong id returned")
	}
}

func TestCreateTeamNoName(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "",
			Owner:               1,
			MembersWithoutOwner: []uint64{},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NAME_TOO_SHORT {
		t.Error("Wrong error")
	}
}

func TestCreateTeamNoOwner(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "test",
			Owner:               0,
			MembersWithoutOwner: []uint64{},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_OWNER_REQUIRED {
		t.Error("Wrong error")
	}
}

func TestCreateTeamNameTooLong(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db), WithMaxTeamNameLength(4))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "testt",
			Owner:               1,
			MembersWithoutOwner: []uint64{},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NAME_TOO_LONG {
		t.Error("Wrong error")
	}
}

func TestCreateTeamNameTooShort(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "t",
			Owner:               1,
			MembersWithoutOwner: []uint64{},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NAME_TOO_SHORT {
		t.Error("Wrong error")
	}
}

func TestCreateTeamTooManyMembers(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db), WithMaxMembers(1))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "test",
			Owner:               1,
			MembersWithoutOwner: []uint64{2},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_TOO_MANY_MEMBERS {
		t.Error("Wrong error")
	}
}

func TestCreateTeamDuplicateMember(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db), WithMaxMembers(2))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "test",
			Owner:               1,
			MembersWithoutOwner: []uint64{2, 2},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NONE {
		t.Error("Wrong error")
	}
}

func TestCreateTeamOwnerInMembers(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db), WithMaxMembers(2))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "test",
			Owner:               1,
			MembersWithoutOwner: []uint64{1, 2},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NONE {
		t.Error("Wrong error")
	}
}

func TestCreateTeamDuplicateMemberOwnerInMembers(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
			return primitive.NilObjectID, nil
		},
	}
	service := NewService(WithDatabase(db), WithMaxMembers(2))
	c := CreateTeamCommand{
		service: service,
		In: &api.CreateTeamRequest{
			Name:                "test",
			Owner:               1,
			MembersWithoutOwner: []uint64{1, 2, 1},
			Data:                map[string]string{},
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
	if c.Out.Error != api.CreateTeamResponse_NONE {
		t.Error("Wrong error")
	}
}
