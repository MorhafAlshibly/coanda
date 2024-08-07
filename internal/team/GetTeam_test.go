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
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

var (
	rankedTeam = []string{"name", "owner", "score", "ranking", "data", "created_at", "updated_at"}
)

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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs("test", 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("test"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
	if c.Out.Team.Owner != 1 {
		t.Fatal("Expected team owner to be 1")
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
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestGetTeamByNameError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs("test", 1).WillReturnError(err)
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("test"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err == nil {
		t.Fatal("Expected error to not be nil")
	}
}

func TestGetTeamByOwner(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 2, 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Owner: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
	if c.Out.Team.Owner != 2 {
		t.Fatal("Expected team owner to be 2")
	}
	if c.Out.Team.Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if !reflect.DeepEqual(c.Out.Team.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamByOwnerNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Owner: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestGetTeamByOwnerError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1).WillReturnError(err)
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Owner: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err == nil {
		t.Fatal("Expected error to not be nil")
	}
}

func TestGetTeamByMember(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1, 1).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 0, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Member: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
	if c.Out.Team.Owner != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if c.Out.Team.Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if !reflect.DeepEqual(c.Out.Team.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamByMemberNotFound(t *testing.T) {
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
		Member: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestGetTeamByMemberError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_team`").WithArgs(2, 1, 1).WillReturnError(err)
	c := NewGetTeamCommand(service, &api.TeamRequest{
		Member: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err == nil {
		t.Fatal("Expected error to not be nil")
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
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
