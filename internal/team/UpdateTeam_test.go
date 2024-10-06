package team

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func TestUpdateTeamNoFieldSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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
		Team: &api.TeamRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamResponse_NO_UPDATE_SPECIFIED {
		t.Fatal("Expected error to be NO_UPDATE_SPECIFIED")
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
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamById(t *testing.T) {
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team` SET `data`=?,`score`=score + ? WHERE (`id` = ?) LIMIT ?")).WithArgs(raw, 2, 5, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Id: conversion.ValueToPointer(uint64(5))},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("UPDATE `team`").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team_with_member`").WithArgs(1, 1, 0).WillReturnError(sql.ErrNoRows)
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Id: conversion.ValueToPointer(uint64(1))},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestUpdateTeamByIdIncrementScoreFalse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team` SET `score`=? WHERE (`id` = ?) LIMIT ?")).WithArgs(2, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Id: conversion.ValueToPointer(uint64(1))},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(false),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamByIdNoScoreSpecified(t *testing.T) {
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team` SET `data`=? WHERE (`id` = ?) LIMIT ?")).WithArgs(raw, 9, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Id: conversion.ValueToPointer(uint64(9))},
		Data:           data,
		IncrementScore: conversion.ValueToPointer(true),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamByIdNoDataSpecified(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team` SET `score`=score + ? WHERE (`id` = ?) LIMIT ?")).WithArgs(2, 9, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:           &api.TeamRequest{Id: conversion.ValueToPointer(uint64(9))},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamByIdNoDataSpecifiedNoIncrementScoreSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTeamCommand(service, &api.UpdateTeamRequest{
		Team:  &api.TeamRequest{Id: conversion.ValueToPointer(uint64(1))},
		Score: conversion.ValueToPointer(int64(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamResponse_INCREMENT_SCORE_NOT_SPECIFIED {
		t.Fatal("Expected error to be INCREMENT_SCORE_NOT_SPECIFIED")
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
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamByMemberId(t *testing.T) {
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
		Team:           &api.TeamRequest{Member: &api.TeamMemberRequest{Id: conversion.ValueToPointer(uint64(1))}},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateTeamByMemberUserId(t *testing.T) {
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
		Team:           &api.TeamRequest{Member: &api.TeamMemberRequest{UserId: conversion.ValueToPointer(uint64(1))}},
		Score:          conversion.ValueToPointer(int64(2)),
		IncrementScore: conversion.ValueToPointer(true),
		Data:           data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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
