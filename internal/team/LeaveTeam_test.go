package team

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func TestLeaveTeamById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `team_member`").WithArgs(3, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(3, 1, 1).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, json.RawMessage("{}"), time.Now(), time.Now()))
	mock.ExpectCommit()
	c := NewLeaveTeamCommand(service, &api.TeamMemberRequest{
		Id: conversion.ValueToPointer(uint64(3)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.LeaveTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestLeaveTeamByIdLastMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `team_member`").WithArgs(3, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(3, 1, 1).WillReturnRows(sqlmock.NewRows(teamMember))
	mock.ExpectExec("DELETE FROM `team`").WithArgs(3, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewLeaveTeamCommand(service, &api.TeamMemberRequest{
		Id: conversion.ValueToPointer(uint64(3)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.LeaveTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestLeaveTeamByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `team_member`").WithArgs(3, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM `team_member`").WithArgs(3, 1, 1).WillReturnRows(sqlmock.NewRows(teamMember).AddRow(1, 2, 3, 1, json.RawMessage("{}"), time.Now(), time.Now()))
	mock.ExpectCommit()
	c := NewLeaveTeamCommand(service, &api.TeamMemberRequest{
		UserId: conversion.ValueToPointer(uint64(3)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.LeaveTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestLeaveTeamByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `team_member`").WithArgs(3, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	c := NewLeaveTeamCommand(service, &api.TeamMemberRequest{
		Id: conversion.ValueToPointer(uint64(3)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.LeaveTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestLeaveTeamByUserIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `team_member`").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	c := NewLeaveTeamCommand(service, &api.TeamMemberRequest{
		UserId: conversion.ValueToPointer(uint64(1)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.LeaveTeamResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestLeaveTeamNoFieldSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewLeaveTeamCommand(service, &api.TeamMemberRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.LeaveTeamResponse_NO_FIELD_SPECIFIED {
		t.Fatal("Expected error to be NO_FIELD_SPECIFIED")
	}
}
