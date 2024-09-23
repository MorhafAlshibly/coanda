package team

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	errorcode "github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestCreateTeamNoScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO team").WithArgs("test", 0, raw).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO team_member").WithArgs(2, 1, 1, raw).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "test",
		FirstMemberUserId: 2,
		Data:              data,
		FirstMemberData:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.CreateTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestCreateTeamFirstMemberExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO team").WithArgs("test", 0, raw).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO team_member").WithArgs(2, 1, 1, raw).WillReturnError(&mysql.MySQLError{Number: errorcode.MySQLErrorCodeDuplicateEntry, Message: "Duplicate entry '2-1' for key 'team_member.team_member_user_id_idx'"})
	mock.ExpectRollback()
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "test",
		FirstMemberUserId: 2,
		Data:              data,
		FirstMemberData:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_FIRST_MEMBER_ALREADY_IN_A_TEAM {
		t.Fatal("Expected error to be FIRST_MEMBER_ALREADY_IN_A_TEAM")
	}
}

func TestCreateTeamNameTaken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO team").WithArgs("test", 0, raw).WillReturnError(&mysql.MySQLError{Number: errorcode.MySQLErrorCodeDuplicateEntry, Message: "Duplicate entry 'test' for key 'team.name'"})
	mock.ExpectRollback()
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "test",
		FirstMemberUserId: 1,
		Data:              data,
		FirstMemberData:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_NAME_TAKEN {
		t.Fatal("Expected error to be NAME_TAKEN")
	}
}

func TestCreateTeamNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "a",
		FirstMemberUserId: 1,
		Data:              nil,
		FirstMemberData:   nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestCreateTeamNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(3), WithMaxTeamNameLength(5))
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "aaaaaaaa",
		FirstMemberUserId: 1,
		Data:              nil,
		FirstMemberData:   nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}

func TestCreateTeamFirstMemberUserIdRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "test",
		FirstMemberUserId: 0,
		Data:              nil,
		FirstMemberData:   nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_FIRST_MEMBER_USER_ID_REQUIRED {
		t.Fatal("Expected error to be FIRST_MEMBER_USER_ID_REQUIRED")
	}
}

func TestCreateTeamDataRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "test",
		FirstMemberUserId: 1,
		Data:              nil,
		FirstMemberData:   nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_DATA_REQUIRED {
		t.Fatal("Expected error to be DATA_REQUIRED")
	}
}

func TestCreateTeamFirstMemberDataRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{
		Name:              "test",
		FirstMemberUserId: 1,
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		},
		FirstMemberData: nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_FIRST_MEMBER_DATA_REQUIRED {
		t.Fatal("Expected error to be FIRST_MEMBER_DATA_REQUIRED")
	}
}

func TestCreateTeamNoInput(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(WithSql(db), WithDatabase(queries))
	c := NewCreateTeamCommand(service, &api.CreateTeamRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTeamResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}
