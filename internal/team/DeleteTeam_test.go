package team

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

func TestDeleteTeamByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM team").WithArgs("test", nil).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Name: conversion.ValueToPointer("test"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestDeleteTeamByOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM team").WithArgs(nil, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Owner: conversion.ValueToPointer(uint64(1)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestDeleteTeamByMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM team").WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTeamCommand(service, &api.TeamRequest{
		Member: conversion.ValueToPointer(uint64(2)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
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
