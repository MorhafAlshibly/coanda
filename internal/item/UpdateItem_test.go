package item

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUpdateItemNoId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateItemResponse_ID_REQUIRED {
		t.Fatal("Expected error to be ID_REQUIRED")
	}
}

func TestUpdateItemNoType(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id: "1",
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateItemResponse_TYPE_REQUIRED {
		t.Fatal("Expected error to be TYPE_REQUIRED")
	}
}

func TestUpdateItemNoUpdateSpecified(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
		},
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateItemResponse_NO_UPDATE_SPECIFIED {
		t.Fatal("Expected error to be NO_UPDATE_SPECIFIED")
	}
}

func TestUpdateItemNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+) AND type = (.+)").WithArgs("1", "type").WillReturnRows(sqlmock.NewRows([]string{}))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
		},
		ExpiresAt: timestamppb.New(time.Now()),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateItemResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestUpdateItemExpiresAt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	inputTime := time.Now().UTC()
	var data json.RawMessage
	mock.ExpectExec("UPDATE").WithArgs(0, data, inputTime, inputTime, "1", "type").WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
		},
		ExpiresAt: timestamppb.New(inputTime),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestUpdateItemData(t *testing.T) {
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
	mock.ExpectExec("UPDATE").WithArgs(1, raw, nil, nil, "1", "type").WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
		},
		Data: data,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestUpdateItemDataAndExpiresAt(t *testing.T) {
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
	inputTime := time.Now().UTC()
	mock.ExpectExec("UPDATE").WithArgs(1, raw, inputTime, inputTime, "1", "type").WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
		},
		Data:      data,
		ExpiresAt: timestamppb.New(inputTime),
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.UpdateItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}
