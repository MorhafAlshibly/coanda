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

func TestGetTournamentUsersNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetTournamentUsersCommand(service, &api.GetTournamentUsersRequest{
		Tournament: conversion.ValueToPointer("t"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_SHORT {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_SHORT")
	}
}

func TestGetTournamentUsersNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTournamentNameLength(5))
	c := NewGetTournamentUsersCommand(service, &api.GetTournamentUsersRequest{
		Tournament: conversion.ValueToPointer("aaaaaaa"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_LONG {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_LONG")
	}
}

func TestGetTournamentUsersNoneFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM ranked_tournament").WillReturnRows(sqlmock.NewRows(rankedTournament))
	c := NewGetTournamentUsersCommand(service, &api.GetTournamentUsersRequest{
		Tournament: conversion.ValueToPointer("test"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTournamentUsersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TournamentUsers) != 0 {
		t.Fatal("Expected users to be empty")
	}
}

func TestGetTournamentUsersTwoUsers(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_tournament").WillReturnRows(sqlmock.NewRows(rankedTournament).AddRow("test", "DAILY", 1, 1, 1, raw, time.Now(), time.Now(), time.Now()).AddRow("test", "DAILY", 2, 1, 1, raw, time.Now(), time.Now(), time.Now()))
	c := NewGetTournamentUsersCommand(service, &api.GetTournamentUsersRequest{
		Tournament: conversion.ValueToPointer("test"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTournamentUsersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TournamentUsers) != 2 {
		t.Fatal("Expected users to have 2 users")
	}
	if c.Out.TournamentUsers[0].UserId != 1 {
		t.Fatal("Expected first user to have id 1")
	}
	if c.Out.TournamentUsers[1].UserId != 2 {
		t.Fatal("Expected second user to have id 2")
	}
}
