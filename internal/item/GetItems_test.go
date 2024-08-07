package item

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

func TestGetItemsNoPagination(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM item WHERE type = (.+)").WithArgs("type", "type", service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(item).AddRow("id", "type", raw, time.Time{}, time.Time{}, time.Time{}))
	c := NewGetItemsCommand(service, &api.GetItemsRequest{
		Type: conversion.ValueToPointer("type"),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Items) != 1 {
		t.Fatal("Expected 1 item")
	}
	if c.Out.Items[0].Id != "id" {
		t.Fatal("Expected id to be id")
	}
	if c.Out.Items[0].Type != "type" {
		t.Fatal("Expected type to be type")
	}
}

func TestGetItemsNoType(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM item WHERE type = (.+)").WithArgs(nil, nil, service.defaultMaxPageLength, 0).WillReturnRows(sqlmock.NewRows(item).AddRow("id", "type", raw, time.Time{}, time.Time{}, time.Time{}))
	c := NewGetItemsCommand(service, &api.GetItemsRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Items) != 1 {
		t.Fatal("Expected 1 item")
	}
}

func TestGetItemsPagination(t *testing.T) {
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
	mock.ExpectQuery("SELECT (.+) FROM item WHERE type = (.+)").WithArgs("type", "type", 1, 1).WillReturnRows(sqlmock.NewRows(item).AddRow("id", "type", raw, time.Time{}, time.Time{}, time.Time{}))
	c := NewGetItemsCommand(service, &api.GetItemsRequest{
		Type: conversion.ValueToPointer("type"),
		Pagination: &api.Pagination{
			Max:  conversion.ValueToPointer(uint32(1)),
			Page: conversion.ValueToPointer(uint64(2)),
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if len(c.Out.Items) != 1 {
		t.Fatal("Expected 1 item")
	}
}
