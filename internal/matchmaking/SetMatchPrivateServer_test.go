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

func Test_SetMatchPrivateServer_EmptyRequest_MatchIdOrMatchmakingTicketRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_SetMatchPrivateServer_EmptyMatchRequest_MatchIdOrMatchmakingTicketRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
		Match: &api.MatchRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_SetMatchPrivateServer_EmptyMatchmakingTicket_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
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
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_SetMatchPrivateServer_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
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
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_SetMatchPrivateServer_NoPrivateServerId_PrivateServerIdRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
		Match: &api.MatchRequest{
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
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_PRIVATE_SERVER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_SetMatchPrivateServer_ByMatchId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(3)),
		},
		PrivateServerId: "123",
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs("123", uint64(3), 1).
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

func Test_SetMatchPrivateServer_MatchDoesntExist_MatchNotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(4)),
		},
		PrivateServerId: "123",
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs("123", uint64(4), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(4), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_SetMatchPrivateServer_PrivateServerAlreadySet_PrivateServerAlreadySetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewSetMatchPrivateServerCommand(service, &api.SetMatchPrivateServerRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(5)),
		},
		PrivateServerId: "123",
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs("123", uint64(5), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(5), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(5), "1234", "PENDING", 1, 1, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(4), 1200, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, 1, json.RawMessage("{}"), time.Now(), time.Now(),
		))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.SetMatchPrivateServerResponse_PRIVATE_SERVER_ALREADY_SET; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
	if got, want := *c.Out.PrivateServerId, "1234"; got != want {
		t.Fatalf("Expected private server id to be %v, got %v", want, got)
	}
}
