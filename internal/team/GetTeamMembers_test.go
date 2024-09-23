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

func TestGetTeamMembersById(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(3, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Id: conversion.ValueToPointer(uint64(3)),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Members) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.Members[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Members[0].UserId != 2 {
		t.Fatal("Expected team user id to be 2")
	}
	if c.Out.Members[0].TeamId != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if !reflect.DeepEqual(c.Out.Members[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersByName(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs("test", 1, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Name: conversion.ValueToPointer("test"),
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Members) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.Members[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Members[0].UserId != 2 {
		t.Fatal("Expected team user id to be 2")
	}
	if c.Out.Members[0].TeamId != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if !reflect.DeepEqual(c.Out.Members[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersByMemberId(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(1, 1, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Member: &api.TeamMemberRequest{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Members) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.Members[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Members[0].UserId != 2 {
		t.Fatal("Expected team user id to be 2")
	}
	if c.Out.Members[0].TeamId != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if !reflect.DeepEqual(c.Out.Members[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersByMemberUserId(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(2, 1, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Member: &api.TeamMemberRequest{
				UserId: conversion.ValueToPointer(uint64(2)),
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Members) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.Members[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Members[0].UserId != 2 {
		t.Fatal("Expected team user id to be 2")
	}
	if c.Out.Members[0].TeamId != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if !reflect.DeepEqual(c.Out.Members[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersMultipleMembers(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(1, 1, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, raw, time.Now(), time.Now()).AddRow(2, 3, 4, 2, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Member: &api.TeamMemberRequest{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Members) != 2 {
		t.Fatal("Expected team members to be 2")
	}
	if c.Out.Members[0].Id != 1 {
		t.Fatal("Expected team id to be 1")
	}
	if c.Out.Members[0].UserId != 2 {
		t.Fatal("Expected team user id to be 2")
	}
	if c.Out.Members[0].TeamId != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if !reflect.DeepEqual(c.Out.Members[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if c.Out.Members[1].Id != 2 {
		t.Fatal("Expected team id to be 2")
	}
	if c.Out.Members[1].UserId != 3 {
		t.Fatal("Expected team user id to be 3")
	}
	if c.Out.Members[1].TeamId != 4 {
		t.Fatal("Expected team id to be 4")
	}
	if !reflect.DeepEqual(c.Out.Members[1].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersNoTeamMembers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(1, 1, service.defaultMaxPageLength).WillReturnRows(sqlmock.NewRows(teamMember))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Member: &api.TeamMemberRequest{
				Id: conversion.ValueToPointer(uint64(1)),
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Members) != 0 {
		t.Fatal("Expected team members to be 0")
	}
}
