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
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs("test", nil, nil, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(teamMember).AddRow("test", 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{Name: conversion.ValueToPointer("test")},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TeamMembers) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.TeamMembers[0].Team != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.TeamMembers[0].UserId != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.TeamMembers[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersByOwner(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs(nil, 1, nil, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(teamMember).AddRow("test", 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{Owner: conversion.ValueToPointer(uint64(1))}})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TeamMembers) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.TeamMembers[0].Team != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.TeamMembers[0].UserId != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.TeamMembers[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMembersByMember(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs(nil, nil, 1, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(teamMember).AddRow("test", 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{Member: conversion.ValueToPointer(uint64(1))}})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TeamMembers) != 1 {
		t.Fatal("Expected team members to be 1")
	}
	if c.Out.TeamMembers[0].Team != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.TeamMembers[0].UserId != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.TeamMembers[0].Data, data) {
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
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs(nil, nil, 1, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(teamMember).AddRow("test", 1, raw, time.Now(), time.Now()).AddRow("test", 2, raw, time.Now(), time.Now()))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{Member: conversion.ValueToPointer(uint64(1))}})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TeamMembers) != 2 {
		t.Fatal("Expected team members to be 2")
	}
	if c.Out.TeamMembers[0].Team != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.TeamMembers[0].UserId != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.TeamMembers[0].Data, data) {
		t.Fatal("Expected team data to be empty")
	}
	if c.Out.TeamMembers[1].Team != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.TeamMembers[1].UserId != 2 {
		t.Fatal("Expected team owner to be 2")
	}
	if !reflect.DeepEqual(c.Out.TeamMembers[1].Data, data) {
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
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs(nil, nil, 1, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(teamMember))
	c := NewGetTeamMembersCommand(service, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{Member: conversion.ValueToPointer(uint64(1))},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMembersResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.TeamMembers) != 0 {
		t.Fatal("Expected team members to be 0")
	}
}
