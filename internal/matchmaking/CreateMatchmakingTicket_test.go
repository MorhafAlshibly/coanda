package matchmaking

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func Test_CreateMatchmakingTicket_NoMatchmakingUsers_MatchmakingUsersRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: nil,
		Arenas:           nil,
		Data:             nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_MATCHMAKING_USERS_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_NoArenas_ArenasRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Arenas: nil,
		Data:   nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_ARENAS_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_NoData_DataRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: nil,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_DATA_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_UserNotFound_UserNotFoundError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_USER_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_MultipleUsersWithOneNotFound_UserNotFoundError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, nil, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(2, 1).WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
			{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_USER_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_UserAlreadyInTicket_UserAlreadyInTicketError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, 1, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_USER_ALREADY_IN_TICKET; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_MultipleUsersWithOneAlreadyInTicket_UserAlreadyInTicketError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, nil, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(2, 1, 2, 1, raw, time.Now(), time.Now()))
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
			{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_USER_ALREADY_IN_TICKET; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_ArenaNotFound_ArenaNotFoundError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, nil, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_ARENA_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_MultipleArenasWithOneNotFound_ArenaNotFoundError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, nil, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).AddRow(1, "test", 2, 3, 5, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").WithArgs(2, 1).WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
			{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_ARENA_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_TooManyPlayersForArena_TooManyPlayersError(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, nil, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(2, nil, 2, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).AddRow(1, "test", 1, 1, 5, raw, time.Now(), time.Now()))
	mock.ExpectRollback()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
			{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_TOO_MANY_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_CreateMatchmakingTicket_ValidInput_Success(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(1, nil, 1, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_user`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(matchmakingUserFields).AddRow(2, nil, 2, 1, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).AddRow(1, "test", 2, 3, 5, raw, time.Now(), time.Now()))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).AddRow(2, "test2", 2, 3, 5, raw, time.Now(), time.Now()))
	mock.ExpectExec("INSERT INTO matchmaking_ticket").WithArgs(raw).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE matchmaking_user").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE matchmaking_user").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectExec("INSERT INTO matchmaking_ticket_arena").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO matchmaking_ticket_arena").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()
	c := NewCreateMatchmakingTicketCommand(service, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: []*api.MatchmakingUserRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
			{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
		Arenas: []*api.ArenaRequest{
			{
				Id: conversion.ValueToPointer(uint64(1)),
			},
			{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := *c.Out.Id, uint64(1); got != want {
		t.Fatalf("Expected id to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.CreateMatchmakingTicketResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
