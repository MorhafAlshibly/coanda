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

func Test_DeleteMatchmakingTicket_EmptyRequest_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatchmakingTicket_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		MatchmakingUser: &api.MatchmakingUserRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchmakingTicketResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatchmakingTicket_ById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(2)),
	})
	mock.ExpectExec("DELETE FROM `matchmaking_ticket`").
		WithArgs(2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchmakingTicketResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatchmakingTicket_ByMatchmakingUserId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		MatchmakingUser: &api.MatchmakingUserRequest{
			Id: conversion.ValueToPointer(uint64(3)),
		},
	})
	mock.ExpectExec("DELETE FROM `matchmaking_ticket`").
		WithArgs(3, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchmakingTicketResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatchmakingTicket_ByMatchmakingUserClientUserId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		MatchmakingUser: &api.MatchmakingUserRequest{
			ClientUserId: conversion.ValueToPointer(uint64(4)),
		},
	})
	mock.ExpectExec("DELETE FROM `matchmaking_ticket`").
		WithArgs(4, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchmakingTicketResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_DeleteMatchmakingTicket_MatchmakingTicketDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteMatchmakingTicketCommand(service, &api.MatchmakingTicketRequest{
		Id: conversion.ValueToPointer(uint64(5)),
	})
	mock.ExpectExec("DELETE FROM `matchmaking_ticket`").
		WithArgs(5, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.DeleteMatchmakingTicketResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
