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

func TestSearchTeamsQueryTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(2))
	c := NewSearchTeamsCommand(service, &api.SearchTeamsRequest{
		Query: "a",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.SearchTeamsResponse_QUERY_TOO_SHORT {
		t.Fatal("Expected error to be QUERY_TOO_SHORT")
	}
}

func TestSearchTeamsQueryTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(2), WithMaxTeamNameLength(6))
	c := NewSearchTeamsCommand(service, &api.SearchTeamsRequest{
		Query: "aaaaaaa",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.SearchTeamsResponse_QUERY_TOO_LONG {
		t.Fatal("Expected error to be QUERY_TOO_LONG")
	}
}

func TestSearchTeamsNoMaxOrPage(t *testing.T) {
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
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(2), WithMaxTeamNameLength(6))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs("aaaa", service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("aaaaaaaa", 1, 10, 1, raw, time.Now(), time.Now()))
	c := NewSearchTeamsCommand(service, &api.SearchTeamsRequest{
		Query: "aaaa",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.SearchTeamsResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected teams to be 1")
	}
	if c.Out.Teams[0].Name != "aaaaaaaa" {
		t.Fatal("Expected team name to be aaaaaaaa")
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

func TestSearchTeamsNoPage(t *testing.T) {
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
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(2), WithMaxTeamNameLength(6))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs("aaaa", 2, 0).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("aaaaaaaa", 1, 10, 1, raw, time.Now(), time.Now()))
	c := NewSearchTeamsCommand(service, &api.SearchTeamsRequest{
		Query:      "aaaa",
		Pagination: &api.Pagination{Max: conversion.ValueToPointer(uint32(2))},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.SearchTeamsResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected teams to be 1")
	}
	if c.Out.Teams[0].Name != "aaaaaaaa" {
		t.Fatal("Expected team name to be aaaaaaaa")
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

func TestSearchTeamsNoMax(t *testing.T) {
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
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(2), WithMaxTeamNameLength(6), WithDefaultMaxPageLength(4))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs("aaaa", 4, 4).WillReturnRows(sqlmock.NewRows(rankedTeam).AddRow("aaaaaaaa", 1, 10, 1, raw, time.Now(), time.Now()))
	c := NewSearchTeamsCommand(service, &api.SearchTeamsRequest{
		Query:      "aaaa",
		Pagination: &api.Pagination{Page: conversion.ValueToPointer(uint64(2))},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.SearchTeamsResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected teams to be 1")
	}
	if c.Out.Teams[0].Name != "aaaaaaaa" {
		t.Fatal("Expected team name to be aaaaaaaa")
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

func TestSearchTeamsNoTeams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(2), WithMaxTeamNameLength(6))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team").WithArgs("aaaa", service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(rankedTeam))
	c := NewSearchTeamsCommand(service, &api.SearchTeamsRequest{
		Query: "aaaa",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.SearchTeamsResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Teams) != 0 {
		t.Fatal("Expected teams to be 0")
	}
}
