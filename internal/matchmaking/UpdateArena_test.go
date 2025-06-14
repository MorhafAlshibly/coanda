package matchmaking

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func Test_UpdateArena_EmptyRequest_ArenaIdOrNameRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_ARENA_ID_OR_NAME_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_EmptyArenaRequest_ArenaIdOrNameRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_ARENA_ID_OR_NAME_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_NameTooShort_NameTooShortError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinArenaNameLength(2))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Name: conversion.ValueToPointer("a"),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_NAME_TOO_SHORT; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_NameTooLong_NameTooLongError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxArenaNameLength(10))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Name: conversion.ValueToPointer("a very long name"),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_NAME_TOO_LONG; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_NoUpdateSpecified_NoUpdateSpecifiedError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_NO_UPDATE_SPECIFIED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_UpdateOnlyOnePlayersValue_IfCapacityChangedMustChangeAllPlayersError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		MinPlayers: conversion.ValueToPointer(uint32(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_IF_CAPACITY_CHANGED_MUST_CHANGE_ALL_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_MinPlayersGreaterThanMaxPlayers_MinPlayersCannotBeGreaterThanMaxPlayersError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		MinPlayers:          conversion.ValueToPointer(uint32(3)),
		MaxPlayers:          conversion.ValueToPointer(uint32(2)),
		MaxPlayersPerTicket: conversion.ValueToPointer(uint32(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_MIN_PLAYERS_CANNOT_BE_GREATER_THAN_MAX_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_MaxPlayersPerTicketLessThanMinPlayers_MaxPlayersPerTicketCannotBeLessThanMinPlayersError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		MinPlayers:          conversion.ValueToPointer(uint32(3)),
		MaxPlayers:          conversion.ValueToPointer(uint32(4)),
		MaxPlayersPerTicket: conversion.ValueToPointer(uint32(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_MAX_PLAYERS_PER_TICKET_CANNOT_BE_LESS_THAN_MIN_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_MaxPlayersPerTicketGreaterThanMaxPlayers_MaxPlayersPerTicketCannotBeGreaterThanMaxPlayersError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		MinPlayers:          conversion.ValueToPointer(uint32(2)),
		MaxPlayers:          conversion.ValueToPointer(uint32(3)),
		MaxPlayersPerTicket: conversion.ValueToPointer(uint32(4)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_MAX_PLAYERS_PER_TICKET_CANNOT_BE_GREATER_THAN_MAX_PLAYERS; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_TicketsInUseWithThisArena_ArenaCurrentlyInUseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		MinPlayers:          conversion.ValueToPointer(uint32(2)),
		MaxPlayers:          conversion.ValueToPointer(uint32(4)),
		MaxPlayersPerTicket: conversion.ValueToPointer(uint32(2)),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(11), uint64(4), "MATCHED", 5, json.RawMessage("{}"), time.Now(), time.Now(),
				4, 4, 1700, 1, json.RawMessage("{}"), time.Now(), time.Now(),
				4, "Arena4", 5, 10, 10, 0, json.RawMessage("{}"), time.Now(), time.Now(),
			))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_ARENA_CURRENTLY_IN_USE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_UpdateDataById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.RawJsonToProtobufStruct(json.RawMessage(`{"key":"value"}`))
	if err != nil {
		t.Fatal(err)
	}
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		Data: data,
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	mock.ExpectExec("UPDATE `matchmaking_arena`").
		WithArgs(json.RawMessage(`{"key":"value"}`), uint64(8), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_UpdateArena_ArenaDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateArenaCommand(service, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id: conversion.ValueToPointer(uint64(999)),
		},
		MinPlayers:          conversion.ValueToPointer(uint32(2)),
		MaxPlayers:          conversion.ValueToPointer(uint32(4)),
		MaxPlayersPerTicket: conversion.ValueToPointer(uint32(2)),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	mock.ExpectExec("UPDATE `matchmaking_arena`").
		WithArgs(uint64(4), uint64(2), uint64(2), uint64(999), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").
		WithArgs(uint64(999), 1).
		WillReturnRows(sqlmock.NewRows(matchmakingArenaFields))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.UpdateArenaResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
