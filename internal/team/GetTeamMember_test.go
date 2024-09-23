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
	teamMember = []string{"id", "user_id", "team_id", "member_number", "data", "joined_at", "updated_at"}
)

func TestGetTeamMemberExists(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMemberCommand(service, &api.TeamMemberRequest{
		UserId: conversion.ValueToPointer(uint64(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMemberResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Member.Id != 1 {
		t.Fatal("Expected team member id to be 1")
	}
	if c.Out.Member.UserId != 2 {
		t.Fatal("Expected team member to be 2")
	}
	if c.Out.Member.TeamId != 3 {
		t.Fatal("Expected team id to be 3")
	}
	if !reflect.DeepEqual(c.Out.Member.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMemberNotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows(teamMember))
	c := NewGetTeamMemberCommand(service, &api.TeamMemberRequest{
		UserId: conversion.ValueToPointer(uint64(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamMemberResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}
