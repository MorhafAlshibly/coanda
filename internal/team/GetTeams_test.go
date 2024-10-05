package team

import (
	"context"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
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
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(1, 0, 1, 0).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember).AddRow(1, "test", 10, 1, raw, time.Now(), time.Now(), 1, 1, 1, raw, time.Now(), time.Now(), 1))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected teams to be 1")
	}
	if c.Out.Teams[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Score != 10 {
		t.Fatal("Expected team score to be 10")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
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
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(1, 0, 2, 0).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember).AddRow(1, "test", 10, 1, raw, time.Now(), time.Now(), 1, 1, 1, raw, time.Now(), time.Now(), 1).AddRow(2, "test2", 5, 2, raw, time.Now(), time.Now(), 2, 2, 1, raw, time.Now(), time.Now(), 1))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Pagination: &api.Pagination{
			Max: conversion.ValueToPointer(uint32(2)),
		}})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 2 {
		t.Fatal("Expected teams to be 2")
	}
	if c.Out.Teams[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Score != 10 {
		t.Fatal("Expected team score to be 10")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if c.Out.Teams[1].Id != 2 {
		t.Fatal("Expected team id to be 2")
	}
	if c.Out.Teams[1].Name != "test2" {
		t.Fatal("Expected team name to be test2")
	}
	if c.Out.Teams[1].Score != 5 {
		t.Fatal("Expected team score to be 5")
	}
	if c.Out.Teams[1].Ranking != 2 {
		t.Fatal("Expected team ranking to be 2")
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
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(1))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(1, 0, 2, 2).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember).AddRow(1, "test", 10, 1, raw, time.Now(), time.Now(), 1, 1, 1, raw, time.Now(), time.Now(), 1).AddRow(2, "test2", 5, 2, raw, time.Now(), time.Now(), 2, 2, 1, raw, time.Now(), time.Now(), 1))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Pagination: &api.Pagination{
			Max:  conversion.ValueToPointer(uint32(2)),
			Page: conversion.ValueToPointer(uint64(2)),
		}})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 2 {
		t.Fatal("Expected teams to be 2")
	}
	if c.Out.Teams[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Score != 10 {
		t.Fatal("Expected team score to be 10")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if c.Out.Teams[1].Id != 2 {
		t.Fatal("Expected team id to be 2")
	}
	if c.Out.Teams[1].Name != "test2" {
		t.Fatal("Expected team name to be test2")
	}
	if c.Out.Teams[1].Score != 5 {
		t.Fatal("Expected team score to be 5")
	}
	if c.Out.Teams[1].Ranking != 2 {
		t.Fatal("Expected team ranking to be 2")
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
		WithSql(db), WithDatabase(queries), WithMaxMaxPageLength(1), WithDefaultMaxPageLength(1))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(1, 0, 1, 0).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember).AddRow(1, "test", 10, 1, raw, time.Now(), time.Now(), 1, 1, 1, raw, time.Now(), time.Now(), 1))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Pagination: &api.Pagination{
			Max: conversion.ValueToPointer(uint32(2)),
		}})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 1 {
		t.Fatal("Expected only 1 team to be returned")
	}
	if c.Out.Teams[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Score != 10 {
		t.Fatal("Expected team score to be 10")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(service.defaultMaxPageLength, 0, service.defaultMaxPageLength, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{
		Pagination: &api.Pagination{
			Page: conversion.ValueToPointer(uint64(2)),
		}})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(service.defaultMaxPageLength, 0, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func Test_GetTeams_MultipleTeamsWithMultipleMembers_TeamsWithMembersReturned(t *testing.T) {
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
		WithSql(db), WithDatabase(queries), WithDefaultMaxPageLength(5))
	mock.ExpectQuery("SELECT (.+) FROM ranked_team_with_member").WithArgs(5, 0, 5, 0).WillReturnRows(sqlmock.NewRows(rankedTeamWithMember).AddRows(
		[]driver.Value{1, "test", 10, 1, raw, time.Now(), time.Now(), 1, 1, 1, raw, time.Now(), time.Now(), 1},
		[]driver.Value{1, "test", 10, 1, raw, time.Now(), time.Now(), 2, 2, 2, raw, time.Now(), time.Now(), 2},
		[]driver.Value{1, "test", 10, 1, raw, time.Now(), time.Now(), 3, 3, 3, raw, time.Now(), time.Now(), 3},
		[]driver.Value{1, "test", 10, 1, raw, time.Now(), time.Now(), 4, 4, 4, raw, time.Now(), time.Now(), 4},
		[]driver.Value{2, "test2", 5, 2, raw, time.Now(), time.Now(), 5, 5, 1, raw, time.Now(), time.Now(), 1},
		[]driver.Value{2, "test2", 5, 2, raw, time.Now(), time.Now(), 6, 6, 2, raw, time.Now(), time.Now(), 2},
		[]driver.Value{2, "test2", 5, 2, raw, time.Now(), time.Now(), 7, 7, 3, raw, time.Now(), time.Now(), 3},
		[]driver.Value{3, "test3", 2, 3, raw, time.Now(), time.Now(), 8, 8, 1, raw, time.Now(), time.Now(), 1},
		[]driver.Value{3, "test3", 2, 3, raw, time.Now(), time.Now(), 9, 9, 2, raw, time.Now(), time.Now(), 2},
		[]driver.Value{3, "test3", 2, 3, raw, time.Now(), time.Now(), 10, 10, 3, raw, time.Now(), time.Now(), 3},
		[]driver.Value{3, "test3", 2, 3, raw, time.Now(), time.Now(), 11, 11, 4, raw, time.Now(), time.Now(), 4},
		[]driver.Value{3, "test3", 2, 3, raw, time.Now(), time.Now(), 12, 12, 5, raw, time.Now(), time.Now(), 5},
		[]driver.Value{4, "test4", 1, 4, raw, time.Now(), time.Now(), 13, 13, 1, raw, time.Now(), time.Now(), 1},
		[]driver.Value{4, "test4", 1, 4, raw, time.Now(), time.Now(), 14, 14, 2, raw, time.Now(), time.Now(), 2},
		[]driver.Value{5, "test5", 0, 5, raw, time.Now(), time.Now(), 15, 15, 1, raw, time.Now(), time.Now(), 1},
		[]driver.Value{5, "test5", 0, 5, raw, time.Now(), time.Now(), 16, 16, 2, raw, time.Now(), time.Now(), 2},
		[]driver.Value{5, "test5", 0, 5, raw, time.Now(), time.Now(), 17, 17, 3, raw, time.Now(), time.Now(), 3},
		[]driver.Value{5, "test5", 0, 5, raw, time.Now(), time.Now(), 18, 18, 4, raw, time.Now(), time.Now(), 4},
		[]driver.Value{5, "test5", 0, 5, raw, time.Now(), time.Now(), 19, 19, 5, raw, time.Now(), time.Now(), 5},
	))
	c := NewGetTeamsCommand(service, &api.GetTeamsRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Teams) != 5 {
		t.Fatal("Expected teams to be 5")
	}
	if c.Out.Teams[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Teams[0].Name != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.Teams[0].Score != 10 {
		t.Fatal("Expected team score to be 10")
	}
	if c.Out.Teams[0].Ranking != 1 {
		t.Fatal("Expected team ranking to be 1")
	}
	if !reflect.DeepEqual(c.Out.Teams[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if len(c.Out.Teams[0].Members) != 4 {
		t.Fatal("Expected 4 members to be returned")
	}
	if c.Out.Teams[0].Members[0].Id != 1 {
		t.Fatal("Expected member id to be 1")
	}
	if c.Out.Teams[0].Members[0].UserId != 1 {
		t.Fatal("Expected member user id to be 1")
	}
	if c.Out.Teams[0].Members[1].Id != 2 {
		t.Fatal("Expected member id to be 2")
	}
	if c.Out.Teams[0].Members[1].UserId != 2 {
		t.Fatal("Expected member user id to be 2")
	}
	if c.Out.Teams[0].Members[2].Id != 3 {
		t.Fatal("Expected member id to be 3")
	}
	if c.Out.Teams[0].Members[2].UserId != 3 {
		t.Fatal("Expected member user id to be 3")
	}
	if c.Out.Teams[0].Members[3].Id != 4 {
		t.Fatal("Expected member id to be 4")
	}
	if c.Out.Teams[0].Members[3].UserId != 4 {
		t.Fatal("Expected member user id to be 4")
	}
	if c.Out.Teams[1].Id != 2 {
		t.Fatal("Expected team id to be 2")
	}
	if c.Out.Teams[1].Name != "test2" {
		t.Fatal("Expected team name to be test2")
	}
	if c.Out.Teams[1].Score != 5 {
		t.Fatal("Expected team score to be 5")
	}
	if c.Out.Teams[1].Ranking != 2 {
		t.Fatal("Expected team ranking to be 2")
	}
	if !reflect.DeepEqual(c.Out.Teams[1].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if len(c.Out.Teams[1].Members) != 3 {
		t.Fatal("Expected 3 members to be returned")
	}
	if c.Out.Teams[1].Members[0].Id != 5 {
		t.Fatal("Expected member id to be 5")
	}
	if c.Out.Teams[1].Members[0].UserId != 5 {
		t.Fatal("Expected member user id to be 5")
	}
	if c.Out.Teams[1].Members[1].Id != 6 {
		t.Fatal("Expected member id to be 6")
	}
	if c.Out.Teams[1].Members[1].UserId != 6 {
		t.Fatal("Expected member user id to be 6")
	}
	if c.Out.Teams[1].Members[2].Id != 7 {
		t.Fatal("Expected member id to be 7")
	}
	if c.Out.Teams[1].Members[2].UserId != 7 {
		t.Fatal("Expected member user id to be 7")
	}
	if c.Out.Teams[2].Id != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if c.Out.Teams[2].Name != "test3" {
		t.Fatal("Expected team name to be test3")
	}
	if c.Out.Teams[2].Score != 2 {
		t.Fatal("Expected team score to be 2")
	}
	if c.Out.Teams[2].Ranking != 3 {
		t.Fatal("Expected team ranking to be 3")
	}
	if !reflect.DeepEqual(c.Out.Teams[2].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if len(c.Out.Teams[2].Members) != 5 {
		t.Fatal("Expected 5 members to be returned")
	}
	if c.Out.Teams[2].Members[0].Id != 8 {
		t.Fatal("Expected member id to be 8")
	}
	if c.Out.Teams[2].Members[0].UserId != 8 {
		t.Fatal("Expected member user id to be 8")
	}
	if c.Out.Teams[2].Members[1].Id != 9 {
		t.Fatal("Expected member id to be 9")
	}
	if c.Out.Teams[2].Members[1].UserId != 9 {
		t.Fatal("Expected member user id to be 9")
	}
	if c.Out.Teams[2].Members[2].Id != 10 {
		t.Fatal("Expected member id to be 10")
	}
	if c.Out.Teams[2].Members[2].UserId != 10 {
		t.Fatal("Expected member user id to be 10")
	}
	if c.Out.Teams[2].Members[3].Id != 11 {
		t.Fatal("Expected member id to be 11")
	}
	if c.Out.Teams[2].Members[3].UserId != 11 {
		t.Fatal("Expected member user id to be 11")
	}
	if c.Out.Teams[2].Members[4].Id != 12 {
		t.Fatal("Expected member id to be 12")
	}
	if c.Out.Teams[2].Members[4].UserId != 12 {
		t.Fatal("Expected member user id to be 12")
	}
	if c.Out.Teams[3].Id != 4 {
		t.Fatal("Expected team id to be 4")
	}
	if c.Out.Teams[3].Name != "test4" {
		t.Fatal("Expected team name to be test4")
	}
	if c.Out.Teams[3].Score != 1 {
		t.Fatal("Expected team score to be 1")
	}
	if c.Out.Teams[3].Ranking != 4 {
		t.Fatal("Expected team ranking to be 4")
	}
	if !reflect.DeepEqual(c.Out.Teams[3].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if len(c.Out.Teams[3].Members) != 2 {
		t.Fatal("Expected 2 members to be returned")
	}
	if c.Out.Teams[3].Members[0].Id != 13 {
		t.Fatal("Expected member id to be 13")
	}
	if c.Out.Teams[3].Members[0].UserId != 13 {
		t.Fatal("Expected member user id to be 13")
	}
	if c.Out.Teams[3].Members[1].Id != 14 {
		t.Fatal("Expected member id to be 14")
	}
	if c.Out.Teams[3].Members[1].UserId != 14 {
		t.Fatal("Expected member user id to be 14")
	}
	if c.Out.Teams[4].Id != 5 {
		t.Fatal("Expected team id to be 5")
	}
	if c.Out.Teams[4].Name != "test5" {
		t.Fatal("Expected team name to be test5")
	}
	if c.Out.Teams[4].Score != 0 {
		t.Fatal("Expected team score to be 0")
	}
	if c.Out.Teams[4].Ranking != 5 {
		t.Fatal("Expected team ranking to be 5")
	}
	if !reflect.DeepEqual(c.Out.Teams[4].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if len(c.Out.Teams[4].Members) != 5 {
		t.Fatal("Expected 5 members to be returned")
	}
	if c.Out.Teams[4].Members[0].Id != 15 {
		t.Fatal("Expected member id to be 15")
	}
	if c.Out.Teams[4].Members[0].UserId != 15 {
		t.Fatal("Expected member user id to be 15")
	}
	if c.Out.Teams[4].Members[1].Id != 16 {
		t.Fatal("Expected member id to be 16")
	}
	if c.Out.Teams[4].Members[1].UserId != 16 {
		t.Fatal("Expected member user id to be 16")
	}
	if c.Out.Teams[4].Members[2].Id != 17 {
		t.Fatal("Expected member id to be 17")
	}
	if c.Out.Teams[4].Members[2].UserId != 17 {
		t.Fatal("Expected member user id to be 17")
	}
	if c.Out.Teams[4].Members[3].Id != 18 {
		t.Fatal("Expected member id to be 18")
	}
	if c.Out.Teams[4].Members[3].UserId != 18 {
		t.Fatal("Expected member user id to be 18")
	}
	if c.Out.Teams[4].Members[4].Id != 19 {
		t.Fatal("Expected member id to be 19")
	}
	if c.Out.Teams[4].Members[4].UserId != 19 {
		t.Fatal("Expected member user id to be 19")
	}
}
