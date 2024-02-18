package tournament

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

func TestUpdateTournamentUserNoTournamentRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED {
		t.Fatal("Expected error to be ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED")
	}
}

func TestUpdateTournamentUserNoTournamentIntervalUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED {
		t.Fatal("Expected error to be ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED")
	}
}

func TestUpdateTournamentUserNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			TournamentIntervalUserId: &api.TournamentIntervalUserId{
				Tournament: "a",
			},
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_SHORT")
	}
}

func TestUpdateTournamentUserNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTournamentNameLength(5))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			TournamentIntervalUserId: &api.TournamentIntervalUserId{
				Tournament: "123456",
			},
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_TOURNAMENT_NAME_TOO_LONG {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_LONG")
	}
}

func TestUpdateTournamentUserNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			TournamentIntervalUserId: &api.TournamentIntervalUserId{
				Tournament: "test",
			},
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestUpdateTournamentUserNoUpdateSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_NO_UPDATE_SPECIFIED {
		t.Fatal("Expected error to be NO_UPDATE_SPECIFIED")
	}
}

func TestUpdateTournamentUserScoreWithoutIncrement(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
		Score:          conversion.ValueToPointer(int64(1)),
		IncrementScore: nil,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_INCREMENT_SCORE_NOT_SPECIFIED {
		t.Fatal("Expected error to be INCREMENT_SCORE_NOT_SPECIFIED")
	}
}

func TestUpdateTournamentUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("UPDATE tournament").WithArgs(raw, uint64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
		Data: data,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}

func TestUpdateTournamentUserScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("UPDATE tournament").WithArgs(1, uint64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
		IncrementScore: conversion.ValueToPointer(false),
		Score:          conversion.ValueToPointer(int64(1)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}

func TestUpdateTournamentUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("UPDATE tournament").WithArgs(1, uint64(1)).WillReturnResult(sqlmock.NewResult(0, 0))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
		Score:          conversion.ValueToPointer(int64(1)),
		IncrementScore: conversion.ValueToPointer(true),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTournamentUserResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestUpdateTournamentUserByTournamentIntervalUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("UPDATE tournament").WithArgs(raw, "test", "DAILY", int64(1), time.Now().Truncate(time.Hour*24).UTC()).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			TournamentIntervalUserId: &api.TournamentIntervalUserId{
				Tournament: "test",
				Interval:   api.TournamentInterval_DAILY,
				UserId:     1,
			},
		},
		Data: data,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}
