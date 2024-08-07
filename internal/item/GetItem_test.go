package item

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

var (
	item = []string{"id", "type", "data", "expires_at", "created_at", "updated_at"}
)

func TestGetItemNoId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetItemCommand(service, &api.ItemRequest{
		Type: "type",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetItemResponse_ID_REQUIRED {
		t.Fatal("Expected error to be ID_REQUIRED")
	}
}

func TestGetItemNoType(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetItemCommand(service, &api.ItemRequest{
		Id: "id",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetItemResponse_TYPE_REQUIRED {
		t.Fatal("Expected error to be TYPE_REQUIRED")
	}
}

func TestGetItemNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetItemCommand(service, &api.ItemRequest{
		Id:   "id",
		Type: "type",
	})
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetItemResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}

func TestGetItemSuccess(t *testing.T) {
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
	c := NewGetItemCommand(service, &api.ItemRequest{
		Id:   "id",
		Type: "type",
	})
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows(item).AddRow("id", "type", raw, time.Time{}, time.Time{}, time.Time{}))
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.Item.Id != "id" {
		t.Fatal("Expected id to be id")
	}
	if c.Out.Item.Type != "type" {
		t.Fatal("Expected type to be type")
	}
}
