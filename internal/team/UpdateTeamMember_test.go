package team

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func TestUpdateTeamMemberNoFieldSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTeamMemberCommand(service, &api.UpdateTeamMemberRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamMemberResponse_NO_FIELD_SPECIFIED {
		t.Fatal("Expected error to be NO_FIELD_SPECIFIED")
	}
}

func TestUpdateTeamMemberNoData(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateTeamMemberCommand(service, &api.UpdateTeamMemberRequest{
		Member: &api.TeamMemberRequest{
			UserId: conversion.ValueToPointer(uint64(1)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamMemberResponse_DATA_REQUIRED {
		t.Fatal("Expected error to be DATA_REQUIRED")
	}
}

func TestUpdateTeamMemberNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team_member` SET `data`=? WHERE (`user_id` = ?) LIMIT ?")).WithArgs(raw, 18, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	c := NewUpdateTeamMemberCommand(service, &api.UpdateTeamMemberRequest{
		Member: &api.TeamMemberRequest{
			UserId: conversion.ValueToPointer(uint64(18)),
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateTeamMemberResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestUpdateTeamMemberById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team_member` SET `data`=? WHERE (`id` = ?) LIMIT ?")).WithArgs(raw, 7, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamMemberCommand(service, &api.UpdateTeamMemberRequest{
		Member: &api.TeamMemberRequest{
			Id: conversion.ValueToPointer(uint64(7)),
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateTeamMemberResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestUpdateTeamMemberByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `team_member` SET `data`=? WHERE (`user_id` = ?) LIMIT ?")).WithArgs(raw, 21, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateTeamMemberCommand(service, &api.UpdateTeamMemberRequest{
		Member: &api.TeamMemberRequest{
			UserId: conversion.ValueToPointer(uint64(21)),
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateTeamMemberResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}
