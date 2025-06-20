package item

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
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
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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

func TestUpdateItemDataRequired(t *testing.T) {
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
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.UpdateItemResponse_DATA_REQUIRED {
		t.Fatal("Expected error to be DATA_REQUIRED")
	}
}

func TestUpdateItemNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+) AND type = (.+)").WithArgs("1", "type").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectRollback()
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
		},
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
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
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(raw, "1", "type").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	c := NewUpdateItemCommand(service, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   "1",
			Type: "type",
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
	if c.Out.Error != api.UpdateItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}
