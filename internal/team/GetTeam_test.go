package team

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

var (
	rankedTeam = []string{"id", "name", "score", "ranking", "data", "created_at", "updated_at"}
)

func TestGetTeamById(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow(1, "test", 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Id: conversion.ValueToPointer(uint64(1)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Team.Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Team.Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if !reflect.DeepEqual(c.Out.Team.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamByName(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs("test", 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow(1, "test", 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("test"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Team.Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Team.Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if !reflect.DeepEqual(c.Out.Team.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamByNameNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs("test", 1).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("test"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestGetTeamByMemberId(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow(1, "test", 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Member: &api.TeamMemberRequest{
			Id: conversion.ValueToPointer(uint64(2)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Team.Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Team.Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if !reflect.DeepEqual(c.Out.Team.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamByMemberIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Member: &api.TeamMemberRequest{
			Id: conversion.ValueToPointer(uint64(2)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestGetTeamByMemberUserId(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow(1, "test", 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Member: &api.TeamMemberRequest{
			UserId: conversion.ValueToPointer(uint64(2)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Team.Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Team.Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if !reflect.DeepEqual(c.Out.Team.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamByMemberUserIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Member: &api.TeamMemberRequest{
			UserId: conversion.ValueToPointer(uint64(2)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestGetTeamNoFieldSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetTeamCommand(service, &api.TeamRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamResponse_NO_FIELD_SPECIFIED {
		t.Fatal(c.Out.Error)
		t.Fatal("Expected error to be NO_FIELD_SPECIFIED")
	}
}

func TestGetTeamNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(3))

	c := NewGetTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("aa"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestGetTeamNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)

	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTeamNameLength(5))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("aaaaaaaa"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}
