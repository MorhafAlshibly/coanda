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

func Test_StartMatch_EmptyRequest_MatchIdOrMatchmakingTicketRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_EmptyMatchRequest_MatchIdOrMatchmakingTicketRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_EmptyMatchmakingTicket_MatchmakingTicketIdOrMatchmakingUserRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
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
	if got, want := c.Out.Error, api.StartMatchResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
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
	if got, want := c.Out.Error, api.StartMatchResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_EmptyStartTime_StartTimeRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
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
	if got, want := c.Out.Error, api.StartMatchResponse_START_TIME_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_StartTimeBeforeCurrentTime_InvalidStartTimeError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(-time.Hour))),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_INVALID_START_TIME; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_StartTimeTooSoon_StartTimeTooSoonError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithStartTimeBuffer(5*time.Minute))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(2 * time.Minute))),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_START_TIME_TOO_SOON; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_LockTimeBeforeCurrentTime_StartTimeTooSoonError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithLockedAtBuffer(5*time.Minute))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(2 * time.Minute))),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_START_TIME_TOO_SOON; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_ByMatchId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithStartTimeBuffer(5*time.Minute), WithLockedAtBuffer(5*time.Minute), WithDefaultMaxPageLength(1))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(10 * time.Minute))),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(8), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(8), "123", "PENDING", 1, 3, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(4), 1200, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, 1, json.RawMessage("{}"), time.Now(), time.Now(),
		))
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(8), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_MatchDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithStartTimeBuffer(5*time.Minute), WithLockedAtBuffer(5*time.Minute), WithDefaultMaxPageLength(1))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(999)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(10 * time.Minute))),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(999), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.StartMatchResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_NotEnoughPlayersToStart_NotEnoughPlayersToStartError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithStartTimeBuffer(5*time.Minute), WithLockedAtBuffer(5*time.Minute), WithDefaultMaxPageLength(1))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(10 * time.Minute))),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(8), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(8), "123", "PENDING", 1, 1, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
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
	if got, want := c.Out.Error, api.StartMatchResponse_NOT_ENOUGH_PLAYERS_TO_START; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_StartedAtAlreadyHasAValue_AlreadyHasStartTimeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithStartTimeBuffer(5*time.Minute), WithLockedAtBuffer(5*time.Minute), WithDefaultMaxPageLength(1))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(10 * time.Minute))),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(8), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(8), "123", "PENDING", 1, 3, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), nil, time.Now(), time.Now(),
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
	if got, want := c.Out.Error, api.StartMatchResponse_ALREADY_HAS_START_TIME; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_StartMatch_PrivateServerNotSet_PrivateServerNotSetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithStartTimeBuffer(5*time.Minute), WithLockedAtBuffer(5*time.Minute), WithDefaultMaxPageLength(1))
	c := NewStartMatchCommand(service, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(8)),
		},
		StartTime: conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(10 * time.Minute))),
	})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(8), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(8), nil, "PENDING", 1, 3, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
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
	if got, want := c.Out.Error, api.StartMatchResponse_PRIVATE_SERVER_NOT_SET; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
