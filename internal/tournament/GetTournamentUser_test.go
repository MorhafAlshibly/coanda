package tournament

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

var (
	rankedTournament = []string{"id", "name", "tournament_interval",
		"user_id", "score", "ranking", "data", "tournament_started_at", "sent_to_third_party_at", "created_at", "updated_at"}
)

func TestGetTournamentUserNoTournamentIntervalUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUserResponse_ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED {
		t.Fatal("Expected error to be ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED")
	}
}

func TestGetTournamentUserTournamentNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{
		TournamentIntervalUserId: &api.TournamentIntervalUserId{
			Tournament: "a",
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_SHORT")
	}
}

func TestGetTournamentUserTournamentNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTournamentNameLength(5))
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{
		TournamentIntervalUserId: &api.TournamentIntervalUserId{
			Tournament: "123456",
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUserResponse_TOURNAMENT_NAME_TOO_LONG {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_LONG")
	}
}

func TestGetTournamentUserNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{
		TournamentIntervalUserId: &api.TournamentIntervalUserId{
			Tournament: "test",
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUserResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestGetTournamentUserById(t *testing.T) {
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
	rows := sqlmock.NewRows(rankedTournament).AddRow(1, "test", "DAILY", 1, 1, 1, raw, time.Now(), time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM `ranked_tournament`").WithArgs(1, 1).WillReturnRows(rows)
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{
		Id: conversion.ValueToPointer(uint64(1)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTournamentUserResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.TournamentUser.Tournament != "test" {
		t.Fatal("Expected tournament to be test")
	}
	if c.Out.TournamentUser.Interval != api.TournamentInterval_DAILY {
		t.Fatal("Expected interval to be DAILY")
	}
	if c.Out.TournamentUser.UserId != 1 {
		t.Fatal("Expected user id to be 1")
	}
	if c.Out.TournamentUser.Score != 1 {
		t.Fatal("Expected score to be 1")
	}
	if c.Out.TournamentUser.Ranking != 1 {
		t.Fatal("Expected ranking to be 1")
	}
	if !reflect.DeepEqual(c.Out.TournamentUser.Data, data) {
		t.Fatal("Expected data to be data")
	}
}

func TestGetTournamentUserByTournamentIntervalUserId(t *testing.T) {
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
	rows := sqlmock.NewRows(rankedTournament).AddRow(1, "test", "DAILY", 1, 1, 1, raw, time.Now(), time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM `ranked_tournament`").WithArgs("test", "DAILY", time.Now().Truncate(24*time.Hour).UTC(), 1, 1).WillReturnRows(rows)
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{
		TournamentIntervalUserId: &api.TournamentIntervalUserId{
			Tournament: "test",
			Interval:   api.TournamentInterval_DAILY,
			UserId:     1,
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTournamentUserResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.TournamentUser.Tournament != "test" {
		t.Fatal("Expected tournament to be test")
	}
	if c.Out.TournamentUser.Interval != api.TournamentInterval_DAILY {
		t.Fatal("Expected interval to be DAILY")
	}
	if c.Out.TournamentUser.UserId != 1 {
		t.Fatal("Expected user id to be 1")
	}
	if c.Out.TournamentUser.Score != 1 {
		t.Fatal("Expected score to be 1")
	}
	if c.Out.TournamentUser.Ranking != 1 {
		t.Fatal("Expected ranking to be 1")
	}
	if !reflect.DeepEqual(c.Out.TournamentUser.Data, data) {
		t.Fatal("Expected data to be data")
	}
}

func TestGetTournamentUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_tournament`").WillReturnError(sql.ErrNoRows)
	c := NewGetTournamentUserCommand(service, &api.TournamentUserRequest{
		TournamentIntervalUserId: &api.TournamentIntervalUserId{
			Tournament: "test",
			UserId:     1,
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUserResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}
