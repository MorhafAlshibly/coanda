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
	matchmakingMatchWithArenaAndTicketFields = []string{
		"match_id", "private_server_id", "match_status", "ticket_count", "user_count",
		"match_data", "locked_at", "started_at", "ended_at", "match_created_at", "match_updated_at",
		"arena_id", "arena_name", "arena_min_players", "arena_max_players_per_ticket", "arena_max_players",
		"arena_data", "arena_created_at", "arena_updated_at",
		"ticket_id", "matchmaking_user_id", "ticket_status", "ticket_user_count", "ticket_number",
		"ticket_data", "expires_at", "ticket_created_at", "ticket_updated_at",
		"client_user_id", "elo", "user_number", "user_data", "user_created_at", "user_updated_at",
		"ticket_arena_id", "ticket_arena_name", "ticket_arena_min_players", "ticket_arena_max_players_per_ticket",
		"ticket_arena_max_players", "arena_number",
		"ticket_arena_data", "ticket_arena_created_at", "ticket_arena_updated_at",
	}
)

func Test_EndMatch_EmptyRequest_MatchIdOrMatchmakingTicketRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.EndMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_EmptyMatchRequest_MatchIdOrMatchmakingTicketRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.EndMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_EmptyTicketRequest_MatchmakingTicketIdOrMatchmakingUserRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
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
	if got, want := c.Out.Error, api.EndMatchResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_EmptyMatchmakingUser_MatchmakingUserIdOrClientUserIdRequired(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
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
	if got, want := c.Out.Error, api.EndMatchResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_NoEndTime_EndTimeRequiredError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.EndMatchResponse_END_TIME_REQUIRED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_TimeInPast_InvalidEndTimeError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	pastTime := conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(-time.Hour)))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
		EndTime: pastTime,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.EndMatchResponse_INVALID_END_TIME; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_ValidInput_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	endTime := conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(time.Hour)))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
		EndTime: endTime,
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs(endTime.AsTime(), uint64(9), endTime.AsTime(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.EndMatchResponse_NONE; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_MatchDoesntExist_NotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	endTime := conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(time.Hour)))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
		EndTime: endTime,
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs(endTime.AsTime(), uint64(9), endTime.AsTime(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(9), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields))
	mock.ExpectRollback()
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, false; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := c.Out.Error, api.EndMatchResponse_NOT_FOUND; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_StartTimeNotSet_StartTimeNotSetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	endTime := conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(time.Hour)))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
		EndTime: endTime,
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs(endTime.AsTime(), uint64(9), endTime.AsTime(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(9), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(9), nil, "PENDING", 1, 1, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now().Add(time.Hour), time.Now(), time.Now(),
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
	if got, want := c.Out.Error, api.EndMatchResponse_START_TIME_NOT_SET; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_MatchAlreadyEnded_AlreadyEndedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	endTime := conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(time.Hour)))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
		EndTime: endTime,
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs(endTime.AsTime(), uint64(9), endTime.AsTime(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(9), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(9), nil, "ENDED", 1, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now().Add(time.Hour), time.Now(), time.Now(),
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
	if got, want := c.Out.Error, api.EndMatchResponse_ALREADY_ENDED; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}

func Test_EndMatch_EndTimeBeforeStartTime_EndTimeBeforeStartTimeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	endTime := conversion.TimeToTimestamppb(conversion.ValueToPointer(time.Now().Add(time.Hour)))
	c := NewEndMatchCommand(service, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: conversion.ValueToPointer(uint64(9)),
		},
		EndTime: endTime,
	})
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `matchmaking_match`").
		WithArgs(endTime.AsTime(), uint64(9), endTime.AsTime(), 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WithArgs(uint64(9), 0, 1, 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).AddRow(
			uint64(9), nil, "STARTED", 1, 1, json.RawMessage("{}"), time.Now().Add(2*time.Hour), time.Now().Add(2*time.Hour), nil, time.Now(), time.Now(),
			uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
			uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now().Add(time.Hour), time.Now(), time.Now(),
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
	if got, want := c.Out.Error, api.EndMatchResponse_END_TIME_BEFORE_START_TIME; got != want {
		t.Fatalf("Expected error to be %v, got %v", want, got)
	}
}
