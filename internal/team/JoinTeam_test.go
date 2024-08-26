package team

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestJoinTeamByName(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs("test", 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	mock.ExpectExec("INSERT INTO team_member").WithArgs("test", 1, raw, "test", service.maxMembers).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Name: conversion.ValueToPointer("test"),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.JoinTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestJoinTeamByOwner(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	mock.ExpectExec("INSERT INTO team_member").WithArgs("test", 1, raw, "test", service.maxMembers).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Owner: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.JoinTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestJoinTeamByMember(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	mock.ExpectExec("INSERT INTO team_member").WithArgs("test", 1, raw, "test", service.maxMembers).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.JoinTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestJoinTeamByMemberNotFound(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.JoinTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestJoinTeamByMemberAlreadyExists(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	mock.ExpectExec("INSERT INTO team_member").WithArgs("test", 1, raw, "test", service.maxMembers).WillReturnError(&mysql.MySQLError{Number: errorcode.MySQLErrorCodeDuplicateEntry, Message: "Duplicate entry '1' for key 'team_member.user_id'"})
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.JoinTeamResponse_ALREADY_IN_A_TEAM {
		t.Fatal("Expected error to be ALREADY_IN_A_TEAM")
	}
}

func TestJoinTeamByMemberTeamFull(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	mock.ExpectExec("INSERT INTO team_member").WithArgs("test", 1, raw, "test", service.maxMembers).WillReturnResult(sqlmock.NewResult(0, 0))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.JoinTeamResponse_TEAM_FULL {
		t.Fatal("Expected error to be TEAM_FULL")
	}
}

func TestJoinTeamByMemberTeamNotFound(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.JoinTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestJoinTeamNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		Data: &structpb.Struct{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.JoinTeamResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestJoinTeamNoData(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewJoinTeamCommand(service, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Member: conversion.ValueToPointer(uint64(1)),
		},
		UserId: 1,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.JoinTeamResponse_DATA_REQUIRED {
		t.Fatal("Expected error to be DATA_REQUIRED")
	}
}
