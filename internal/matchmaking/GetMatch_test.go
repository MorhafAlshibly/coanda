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

func Test_GetMatch_EmptyRequest_MatchIdOrMatchmakingTicketRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchCommand(service, &api.GetMatchRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetMatch_EmptyMatchRequest_MatchIdOrMatchmakingTicketRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchCommand(service, &api.GetMatchRequest{
		Match: &api.MatchRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetMatch_EmptyMatchmakingTicketRequest_MatchmakingTicketIdOrMatchmakingUserRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchCommand(service, &api.GetMatchRequest{
		Match: &api.MatchRequest{
			MatchmakingTicket: &api.MatchmakingTicketRequest{},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetMatchResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetMatch_EmptyMatchmakingUserRequest_MatchmakingUserIdOrClientUserIdRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchCommand(service, &api.GetMatchRequest{
		Match: &api.MatchRequest{
			MatchmakingTicket: &api.MatchmakingTicketRequest{
				MatchmakingUser: &api.MatchmakingUserRequest{},
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetMatchResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_GetMatch_ByMatchId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewGetMatchCommand(service, &api.GetMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
	})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(9), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(9), nil, "PENDING", 1, 1, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(4), 1200, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, 1, json.RawMessage("{}"), time.Now(), time.Now(),
		))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Match.Id, uint64(9); got != want {
		t.Fatalf("Expected match id to be %d, got %d", want, got)
	}
}

func Test_GetMatch_MatchDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewGetMatchCommand(service, &api.GetMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(999)),
		},
	})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(999), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.GetMatchResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
