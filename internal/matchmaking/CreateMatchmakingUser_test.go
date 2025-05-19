package matchmaking

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/go-sql-driver/mysql"
)

var (
	matchmakingUserFields  = []string{"id", "matchmaking_ticket_id", "client_user_id", "elo", "data", "created_at", "updated_at"}
	matchmakingArenaFields = []string{"id", "name", "min_players", "max_players_per_ticket", "max_players", "data", "created_at", "updated_at"}
)

func Test_CreateMatchmakingUser_NoClientUserId_ClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingUserCommand(service, &api.CreateMatchmakingUserRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingUserResponse_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingUser_NoData_DataRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingUserCommand(service, &api.CreateMatchmakingUserRequest{
		ClientUserId: 1,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingUserResponse_DATA_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingUser_UserAlreadyExists_AlreadyExistsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	raw := json.RawMessage("{}")
	data, err := conversion.RawJsonToProtobufStruct(raw)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingUserCommand(service, &api.CreateMatchmakingUserRequest{
		ClientUserId: 1,
		Data:         data,
	})
	mock.ExpectExec("INSERT INTO matchmaking_user").
		WithArgs(1, 0, raw).
		WillReturnError(&mysql.MySQLError{
			Number:  errorcode.MySQLErrorCodeDuplicateEntry,
			Message: "Duplicate entry '1' for key 'matchmaking_user.client_user_id'",
		})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingUserResponse_ALREADY_EXISTS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingUser_ValidInput_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	raw := json.RawMessage("{}")
	data, err := conversion.RawJsonToProtobufStruct(raw)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingUserCommand(service, &api.CreateMatchmakingUserRequest{
		ClientUserId: 1,
		Data:         data,
	})
	mock.ExpectExec("INSERT INTO matchmaking_user").
		WithArgs(1, 0, raw).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := *c.Out.Id, uint64(1); got != want {
		t.Fatalf("Expected ID to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingUserResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
