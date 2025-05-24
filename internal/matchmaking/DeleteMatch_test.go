package matchmaking

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

var (
	matchmakingMatchFields = []string{"match_id", "private_server_id", "match_status", "ticket_count", "user_count", "match_data", "locked_at", "started_at", "ended_at", "match_created_at", "match_updated_at"}
)

func Test_DeleteMatch_EmptyRequest_MatchIdOrMatchmakingTicketRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_EmptyTicketRequest_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
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
	if got, want := c.Out.Error, api.DeleteMatchResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_ByMatchId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
		Id: conversion.ValueToPointer(uint64(1)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `matchmaking_match`").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_ByMatchmakingTicketId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: conversion.ValueToPointer(uint64(3)),
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `matchmaking_match`").
		WithArgs(3, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_ByMatchmakingUserId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			MatchmakingUser: &api.MatchmakingUserRequest{
				Id: conversion.ValueToPointer(uint64(2)),
			},
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `matchmaking_match`").
		WithArgs(2, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_ByMatchmakingClientUserId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			MatchmakingUser: &api.MatchmakingUserRequest{
				ClientUserId: conversion.ValueToPointer(uint64(4)),
			},
		},
	})
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `matchmaking_match`").
		WithArgs(4, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatch_MatchDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchCommand(service, &api.MatchRequest{
		Id: conversion.ValueToPointer(uint64(10)),
	})
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `matchmaking_match`").
		WithArgs(10, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(10, 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchFields))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
