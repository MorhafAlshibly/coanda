package record

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

var (
	rankedRecord = []string{"id", "name", "user_id", "record", "ranking", "data", "created_at", "updated_at"}
)

func TestGetRecordNoFieldSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetRecordCommand(service, &api.RecordRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordResponse_ID_OR_NAME_USER_ID_REQUIRED {
		t.Fatal("Expected error to be ID_OR_NAME_USER_ID_REQUIRED")
	}
}

func TestGetRecordNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetRecordCommand(service, &api.RecordRequest{
		NameUserId: &api.NameUserId{
			Name: "t",
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestGetRecordNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxRecordNameLength(5))
	c := NewGetRecordCommand(service, &api.RecordRequest{
		NameUserId: &api.NameUserId{
			Name: "aaaaaaa",
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}

func TestGetRecordNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetRecordCommand(service, &api.RecordRequest{
		NameUserId: &api.NameUserId{
			Name: "test",
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestGetRecordSuccess(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_record`").
		WithArgs("test", 1, 1).
		WillReturnRows(sqlmock.NewRows(rankedRecord).AddRow(1, "test", 1, 1, 1, raw, time.Time{}, time.Time{}))
	c := NewGetRecordCommand(service, &api.RecordRequest{
		NameUserId: &api.NameUserId{
			Name:   "test",
			UserId: 1,
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetRecordResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Record.Name != "test" {
		t.Fatal("Expected name to be test")
	}
	if c.Out.Record.Record != 1 {
		t.Fatal("Expected record to be 1")
	}
	if c.Out.Record.UserId != 1 {
		t.Fatal("Expected user_id to be 1")
	}
}

func TestGetRecordById(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM `ranked_record`").
		WithArgs(uint64(1), 1).
		WillReturnRows(sqlmock.NewRows(rankedRecord).AddRow(1, "test", 1, 1, 1, raw, time.Time{}, time.Time{}))
	c := NewGetRecordCommand(service, &api.RecordRequest{
		Id: conversion.ValueToPointer(uint64(1)),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetRecordResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Record.Name != "test" {
		t.Fatal("Expected name to be test")
	}
	if c.Out.Record.Record != 1 {
		t.Fatal("Expected record to be 1")
	}
	if c.Out.Record.UserId != 1 {
		t.Fatal("Expected user_id to be 1")
	}
}

func TestGetRecordNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM `ranked_record`").
		WithArgs("test", 1, 1).WillReturnError(sql.ErrNoRows)
	c := NewGetRecordCommand(service, &api.RecordRequest{
		NameUserId: &api.NameUserId{
			Name:   "test",
			UserId: 1,
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}
