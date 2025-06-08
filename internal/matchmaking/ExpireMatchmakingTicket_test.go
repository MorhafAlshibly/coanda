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

var (
	matchmakingTicketFields = []string{
		"ticket_id", "matchmaking_match_id", "status", "user_count", "ticket_data",
		"expires_at", "ticket_created_at", "ticket_updated_at",
		"matchmaking_user_id", "client_user_id", "elo", "user_number", "user_data",
		"user_created_at", "user_updated_at",
		"arena_id", "arena_name", "arena_min_players", "arena_max_players_per_ticket",
		"arena_max_players", "arena_number", "arena_data", "arena_created_at", "arena_updated_at",
	}
)

func Test_ExpireMatchmakingTicket_EmptyRequest_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.ExpireMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_ExpireMatchmakingTicket_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		MatchmakingUser: &api.MatchmakingUserRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.ExpireMatchmakingTicketResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_ExpireMatchmakingTicket_ById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(6)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(uint64(6), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
}

func Test_ExpireMatchmakingTicket_TicketDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(7)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(uint64(7), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(7), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.ExpireMatchmakingTicketResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_ExpireMatchmakingTicket_TicketAlreadyExpired_AlreadyExpiredError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(8)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(uint64(8), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(8), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(8), nil, "EXPIRED", 2, json.RawMessage("{}"),
				time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				1, 1, 0, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				1, "Arena1", 2, 4, 4, 0, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
			))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.ExpireMatchmakingTicketResponse_ALREADY_EXPIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_ExpireMatchmakingTicket_TicketAlreadyMatched_AlreadyMatchedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(9)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(uint64(9), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(9), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(9), uint64(2), "MATCHED", 3, json.RawMessage("{}"),
				time.Now().Add(time.Hour), time.Now(), time.Now(),
				2, 2, 1500, 1, json.RawMessage("{}"), time.Now(), time.Now(),
				2, "Arena2", 3, 6, 6, 0, json.RawMessage("{}"), time.Now(), time.Now(),
			))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.ExpireMatchmakingTicketResponse_ALREADY_MATCHED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_ExpireMatchmakingTicket_TicketAlreadyEnded_AlreadyEndedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewExpireMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(10)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(uint64(10), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(10), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(10), uint64(3), "ENDED", 4, json.RawMessage("{}"),
				time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, 3, 1600, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, "Arena3", 4, 8, 8, 0, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
			))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.ExpireMatchmakingTicketResponse_ALREADY_ENDED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
