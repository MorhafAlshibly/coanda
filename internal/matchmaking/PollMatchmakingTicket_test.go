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

func Test_PollMatchmakingTicket_EmptyRequest_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_EmptyMatchmakingTicketRequest_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			MatchmakingUser: &api.MatchmakingUserRequest{},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_ById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: conversion.ValueToPointer(uint64(6)),
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(sqlmock.AnyArg(), uint64(6), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(6), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(6), nil, "PENDING", 4, json.RawMessage("{}"), time.Now().Add(time.Hour), time.Now(), time.Now(),
				3, 3, 1600, 1, json.RawMessage("{}"), time.Now(), time.Now(),
				3, "Arena3", 4, 8, 8, 0, json.RawMessage("{}"), time.Now(), time.Now(),
			))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.MatchmakingTicket.Id, uint64(6); got != want {
		t.Fatalf("Expected MatchmakingTicket ID to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_TicketDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: conversion.ValueToPointer(uint64(7)),
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(sqlmock.AnyArg(), uint64(7), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(7), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_AlreadyExpired_AlreadyExpiredError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(sqlmock.AnyArg(), uint64(8), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(8), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(8), nil, "EXPIRED", 2, json.RawMessage("{}"),
				time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, 3, 1600, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, "Arena3", 4, 8, 8, 0, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
			))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_ALREADY_EXPIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_AlreadyMatched_AlreadyMatchedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(sqlmock.AnyArg(), uint64(9), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(9), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(9), uint64(3), "MATCHED", 2, json.RawMessage("{}"),
				time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, 3, 1600, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, "Arena3", 4, 8, 8, 0, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
			))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_ALREADY_MATCHED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_PollMatchmakingTicket_AlreadyEnded_AlreadyEndedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewPollMatchmakingTicketCommand(service, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: conversion.ValueToPointer(uint64(10)),
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_ticket`").
		WithArgs(sqlmock.AnyArg(), uint64(10), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(10), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(10), uint64(3), "ENDED", 2, json.RawMessage("{}"),
				time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, 3, 1600, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, "Arena3", 4, 8, 8, 0, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
			))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.PollMatchmakingTicketResponse_ALREADY_ENDED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
