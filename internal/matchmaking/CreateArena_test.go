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

func Test_CreateArena_NameTooShort_NameTooShortError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinArenaNameLength(2))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name: "t",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_NAME_TOO_SHORT; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_NameTooLong_NameTooLongError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxArenaNameLength(10))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name: "this is a very long name that exceeds the maximum length",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_NAME_TOO_LONG; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_NoMinPlayers_MinPlayersRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name: "test",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_MIN_PLAYERS_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_NoMaxPlayersPerTicket_MaxPlayersPerTicketRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:       "test",
		MinPlayers: 5,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_MAX_PLAYERS_PER_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_NoMaxPlayers_MaxPlayersRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:                "test",
		MinPlayers:          5,
		MaxPlayersPerTicket: 2,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_MAX_PLAYERS_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_MinPlayersGreaterThanMaxPlayers_MinPlayersCannotBeGreaterThanMaxPlayersError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:                "test",
		MinPlayers:          10,
		MaxPlayersPerTicket: 2,
		MaxPlayers:          5,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_MIN_PLAYERS_CANNOT_BE_GREATER_THAN_MAX_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_MaxPlayersPerTicketGreaterThanMaxPlayers_MaxPlayersPerTicketCannotBeGreaterThanMaxPlayersError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:                "test",
		MinPlayers:          5,
		MaxPlayersPerTicket: 10,
		MaxPlayers:          5,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_MAX_PLAYERS_PER_TICKET_CANNOT_BE_GREATER_THAN_MAX_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_NoData_DataRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:                "test",
		MinPlayers:          2,
		MaxPlayersPerTicket: 3,
		MaxPlayers:          5,
		Data:                nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_DATA_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_ValidInput_Success(t *testing.T) {
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
	mock.ExpectExec("INSERT INTO matchmaking_arena").WithArgs("test", 2, 3, 5, raw).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:                "test",
		MinPlayers:          2,
		MaxPlayersPerTicket: 3,
		MaxPlayers:          5,
		Data:                data,
	})
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
	if got, want := c.Out.Error, api.CreateArenaResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateArena_DuplicateName_AlreadyExistsError(t *testing.T) {
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
	mock.ExpectExec("INSERT INTO matchmaking_arena").WithArgs("test", 2, 3, 5, raw).WillReturnError(&mysql.MySQLError{
		Number:  uint16(errorcode.MySQLErrorCodeDuplicateEntry),
		Message: "Duplicate entry 'test' for key 'matchmaking_arena.name'",
	})
	c := NewCreateArenaCommand(service, &api.CreateArenaRequest{
		Name:                "test",
		MinPlayers:          2,
		MaxPlayersPerTicket: 3,
		MaxPlayers:          5,
		Data:                data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateArenaResponse_ALREADY_EXISTS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
