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

func Test_GetArena_EmptyRequest_ArenaIdOrNameRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetArenaCommand(service, &api.ArenaRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetArenaResponse_ARENA_ID_OR_NAME_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetArena_NameTooShort_NameTooShortError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinArenaNameLength(2))
	c := NewGetArenaCommand(service, &api.ArenaRequest{
		Name: conversion.ValueToPointer("a"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetArenaResponse_NAME_TOO_SHORT; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetArena_NameTooLong_NameTooLongError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxArenaNameLength(10))
	c := NewGetArenaCommand(service, &api.ArenaRequest{
		Name: conversion.ValueToPointer("a very long name"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetArenaResponse_NAME_TOO_LONG; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetArena_ById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetArenaCommand(service, &api.ArenaRequest{
		Id: conversion.ValueToPointer(uint64(6)),
	})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").
		WithArgs(uint64(6), 1).
		WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).
			AddRow(uint64(6), "Test Arena", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now()))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetArenaResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
	if got, want := c.Out.Arena.Id, uint64(6); got != want {
		t.Fatalf("Expected arena id to be %d, got %d", want, got)
	}
	if got, want := c.Out.Arena.Name, "Test Arena"; got != want {
		t.Fatalf("Expected arena name to be %s, got %s", want, got)
	}
}

func Test_GetArena_ByName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetArenaCommand(service, &api.ArenaRequest{
		Name: conversion.ValueToPointer("Test Arena"),
	})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").
		WithArgs("Test Arena", 1).
		WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).
			AddRow(uint64(6), "Test Arena", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now()))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetArenaResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
	if got, want := c.Out.Arena.Id, uint64(6); got != want {
		t.Fatalf("Expected arena id to be %d, got %d", want, got)
	}
	if got, want := c.Out.Arena.Name, "Test Arena"; got != want {
		t.Fatalf("Expected arena name to be %s, got %s", want, got)
	}
}

func Test_GetArena_ArenaDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetArenaCommand(service, &api.ArenaRequest{
		Id: conversion.ValueToPointer(uint64(7)),
	})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_arena`").
		WithArgs(uint64(7), 1).
		WillReturnError(sql.ErrNoRows)
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetArenaResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
