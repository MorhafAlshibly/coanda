package team

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

func TestUpdateTeamNoTeam(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamResponse_NO_FIELD_SPECIFIED {
		t.Fatal("Expected error to be NO_FIELD_SPECIFIED")
	}
}

func TestUpdateTeamNoUpdateSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team: &api.TeamRequest{},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamResponse_NO_FIELD_SPECIFIED {
		t.Fatal("Expected error to be NO_FIELD_SPECIFIED")
	}
}

func TestUpdateTeamNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(3))

	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Name: conversion.ValueToPointer("aa")},
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
	if c.Out.Error != api.UpdateTeamResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestUpdateTeamNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(3), WithMaxTeamNameLength(5))

	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Name: conversion.ValueToPointer("aaaaaaaa")},
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
	if c.Out.Error != api.UpdateTeamResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}

func TestUpdateTeamByName(t *testing.T) {
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
	mock.ExpectExec("UPDATE `team`").WithArgs(raw, 2, "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Name: conversion.ValueToPointer("test")},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestUpdateTeamByOwner(t *testing.T) {
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
	mock.ExpectExec("UPDATE `team`").WithArgs(raw, 2, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Owner: conversion.ValueToPointer(uint64(1))},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestUpdateTeamByMember(t *testing.T) {
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
	mock.ExpectExec("UPDATE `team`").WithArgs(raw, 2, 1, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Member: conversion.ValueToPointer(uint64(1))},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}
