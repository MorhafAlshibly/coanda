package tournament

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/go-sql-driver/mysql"
)

func TestCreateTournamentUserTournamentNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTournamentUserCommand(service, &api.CreateTournamentUserRequest{
		Tournament: "t",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_SHORT")
	}
}

func TestCreateTournamentUserTournamentNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTournamentNameLength(5))
	c := NewCreateTournamentUserCommand(service, &api.CreateTournamentUserRequest{
		Tournament: "aaaaaaa",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTournamentUserResponse_TOURNAMENT_NAME_TOO_LONG {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_LONG")
	}
}

func TestCreateTournamentUserNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTournamentUserCommand(service, &api.CreateTournamentUserRequest{
		Tournament: "test",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTournamentUserResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestCreateTournamentUserDataRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateTournamentUserCommand(service, &api.CreateTournamentUserRequest{
		Tournament: "test",
		UserId:     1,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTournamentUserResponse_DATA_REQUIRED {
		t.Fatal("Expected error to be DATA_REQUIRED")
	}
}

func TestCreateTournamentUserSuccess(t *testing.T) {
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
	mock.ExpectExec("INSERT INTO tournament").WithArgs("test", "DAILY", int64(1), int64(0), raw, service.getTournamentStartDate(time.Now(), api.TournamentInterval_DAILY)).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewCreateTournamentUserCommand(service, &api.CreateTournamentUserRequest{
		Tournament: "test",
		Interval:   api.TournamentInterval_DAILY,
		UserId:     1,
		Data:       data,
		Score:      conversion.ValueToPointer(int64(0)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.CreateTournamentUserResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestCreateTournamentUserAlreadyExists(t *testing.T) {
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
	mock.ExpectExec("INSERT INTO tournament").WithArgs("test", "DAILY", int64(1), int64(0), raw, service.getTournamentStartDate(time.Now(), api.TournamentInterval_DAILY)).WillReturnError(&mysql.MySQLError{Number: errorcodes.MySQLErrorCodeDuplicateEntry})
	c := NewCreateTournamentUserCommand(service, &api.CreateTournamentUserRequest{
		Tournament: "test",
		Interval:   api.TournamentInterval_DAILY,
		UserId:     1,
		Data:       data,
		Score:      conversion.ValueToPointer(int64(0)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateTournamentUserResponse_ALREADY_EXISTS {
		t.Fatal("Expected error to be ALREADY_EXISTS")
	}
}
