package record

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

func TestGetRecordsNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetRecordsCommand(service, &api.GetRecordsRequest{
		Name: conversion.ValueToPointer("t"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordsResponse_NAME_TOO_SHORT {
		t.Fatal("Expected error to be NAME_TOO_SHORT")
	}
}

func TestGetRecordsNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxRecordNameLength(5))
	c := NewGetRecordsCommand(service, &api.GetRecordsRequest{
		Name: conversion.ValueToPointer("aaaaaaa"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetRecordsResponse_NAME_TOO_LONG {
		t.Fatal("Expected error to be NAME_TOO_LONG")
	}
}

func TestGetRecordsSuccess(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM ranked_record").WillReturnRows(sqlmock.NewRows(rankedRecord).AddRow("test", 1, 1, 1, raw, time.Time{}, time.Time{}))
	c := NewGetRecordsCommand(service, &api.GetRecordsRequest{
		Name: conversion.ValueToPointer("name"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetRecordsResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if len(c.Out.Records) != 1 {
		t.Fatal("Expected records to have length 1")
	}
	if c.Out.Records[0].Name != "test" {
		t.Fatal("Expected name to be test")
	}
	if c.Out.Records[0].UserId != 1 {
		t.Fatal("Expected user id to be 1")
	}
	if c.Out.Records[0].Record != 1 {
		t.Fatal("Expected record to be 1")
	}
	if !reflect.DeepEqual(c.Out.Records[0].Data, data) {
		t.Fatal("Expected data to be equal")
	}
}
