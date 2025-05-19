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
