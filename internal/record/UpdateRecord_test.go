package record

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func TestUpdateRecordNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name: "t",
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateRecordResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestUpdateRecordNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxRecordNameLength(5))
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name: "aaaaaaa",
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateRecordResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}

func TestUpdateRecordEmptyRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateRecordResponse_ID_OR_NAME_USER_ID_REQUIRED {
		t.Fatal("Expected error to be ID_OR_NAME_USER_ID_REQUIRED")
	}
}

func TestUpdateRecordNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name: "test",
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateRecordResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestUpdateRecordNoRecordOrData(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name:   "test",
				UserId: 1,
			},
		},
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateRecordResponse_NO_UPDATE_SPECIFIED {
		t.Fatal("Expected error to be NO_UPDATE_SPECIFIED")
	}
}

func TestUpdateRecordNoRecord(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `record`").WithArgs(raw, "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name:   "test",
				UserId: 1,
			},
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}

func TestUpdateRecordNoData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `record`").WithArgs(2, "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name:   "test",
				UserId: 1,
			},
		},
		Record: conversion.ValueToPointer(uint64(2)),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}

func TestUpdateRecordRecordAndData(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `record`").WithArgs(raw, 2, "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			NameUserId: &api.NameUserId{
				Name:   "test",
				UserId: 1,
			},
		},
		Record: conversion.ValueToPointer(uint64(2)),
		Data:   data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}

func TestUpdateRecordById(t *testing.T) {
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `record`").WithArgs(raw, uint64(1), "", 0).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateRecordCommand(service, &api.UpdateRecordRequest{
		Request: &api.RecordRequest{
			Id: conversion.ValueToPointer(uint64(1)),
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}
