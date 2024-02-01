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

func TestGetTeamsDefaultSettings(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs(service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected teams to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if c.Out.Teams[0].Owner != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamsMultipleTeams(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs(2, 0).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()).AddRow("test2", 2, 5, 2, raw, time.Now(), time.Now()))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Max: conversion.ValueToPointer(uint32(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 2 {
		t.Fatal("Expected teams to be 2")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if c.Out.Teams[0].Owner != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if c.Out.Teams[1].Name != "test2" {
		t.Fatal("Expected team name to be test2")
	}
	if c.Out.Teams[1].Ranking != 2 {
		t.Fatal("Expected team ranking to be 2")
	}
	if c.Out.Teams[1].Owner != 2 {
		t.Fatal("Expected team owner to be 2")
	}
	if !reflect.DeepEqual(c.Out.Teams[1].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamsMultipleTeamsWithPage(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs(2, 2).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()).AddRow("test2", 2, 5, 2, raw, time.Now(), time.Now()))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Max:  conversion.ValueToPointer(uint32(2)),
		Page: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 2 {
		t.Fatal("Expected teams to be 2")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if c.Out.Teams[0].Owner != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if c.Out.Teams[1].Name != "test2" {
		t.Fatal("Expected team name to be test2")
	}
	if c.Out.Teams[1].Ranking != 2 {
		t.Fatal("Expected team ranking to be 2")
	}
	if c.Out.Teams[1].Owner != 2 {
		t.Fatal("Expected team owner to be 2")
	}
	if !reflect.DeepEqual(c.Out.Teams[1].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamsMultipleTeamsWithTooLargeMax(t *testing.T) {
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
		WithSql(db), WithDatabase(queries), WithMaxMaxPageLength(1))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs(1, 0).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("test", 1, 10, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Max: conversion.ValueToPointer(uint32(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected teams to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if c.Out.Teams[0].Owner != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamsMultipleTeamsWithTooLargePage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs(service.defaultMaxPageLength, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Page: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 0 {
		t.Fatal("Expected teams to be 0")
	}
}

func TestGetTeamsNoTeams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs(service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 0 {
		t.Fatal("Expected teams to be 0")
	}
}
