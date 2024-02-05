package tournament

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

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
			Tournament: "t",
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
			Tournament: "aaaaaaa",
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
			Tournament: "test",
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

func TestUpdateTournamentUserEmptyTournamentRequest(t *testing.T) {
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
	if c.Out.Error != api.UpdateTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_SHORT")
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
			Tournament: "test",
			Interval:   0,
			UserId:     1,
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
			Tournament: "test",
			Interval:   0,
			UserId:     1,
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tournament").WithArgs(raw, "test", "DAILY", int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Tournament: "test",
			Interval:   0,
			UserId:     1,
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tournament").WithArgs(1, 1, 0, "test", "DAILY", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Tournament: "test",
			Interval:   0,
			UserId:     1,
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tournament").WithArgs(1, 1, 1, "test", "DAILY", 1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	c := NewUpdateTournamentUserCommand(service, &api.UpdateTournamentUserRequest{
		Tournament: &api.TournamentUserRequest{
			Tournament: "test",
			Interval:   0,
			UserId:     1,
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
