package team

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func TestDeleteTeamNoNameMember(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TeamResponse_NO_FIELD_SPECIFIED {
		t.Fatal("Expected error to be NO_FIELD_SPECIFIED")
	}
}

func TestDeleteTeamNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMinTeamNameLength(3))

	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("aa"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TeamResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestDeleteTeamNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTeamNameLength(5))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("aaaaaaa"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TeamResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}

func TestDeleteTeamById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM `team`").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Id: conversion.ValueToPointer(uint64(1)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.TeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestDeleteTeamByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM `team`").WithArgs("test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("test"),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.TeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestDeleteTeamByMemberId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM `team`").WithArgs(2, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
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
	if c.Out.Error != api.TeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestDeleteTeamByMemberUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM `team`").WithArgs(2, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
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
	if c.Out.Error != api.TeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}
