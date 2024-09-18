package item

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateItemNoId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateItemCommand(service, &api.CreateItemRequest{
		Type: "type",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateItemResponse_ID_REQUIRED {
		t.Fatal("Expected error to be ID_REQUIRED")
	}
}

func TestCreateItemNoType(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateItemCommand(service, &api.CreateItemRequest{
		Id: "id",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateItemResponse_TYPE_REQUIRED {
		t.Fatal("Expected error to be TYPE_REQUIRED")
	}
}

func TestCreateItemNoData(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewCreateItemCommand(service, &api.CreateItemRequest{
		Id:   "id",
		Type: "type",
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateItemResponse_DATA_REQUIRED {
		t.Fatal("Expected error to be DATA_REQUIRED")
	}
}

func TestCreateItemAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	c := NewCreateItemCommand(service, &api.CreateItemRequest{
		Id:   "id",
		Type: "type",
		Data: data,
	})
	mock.ExpectExec("INSERT INTO item").WillReturnError(&mysql.MySQLError{Number: errorcode.MySQLErrorCodeDuplicateEntry})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.CreateItemResponse_ALREADY_EXISTS {
		t.Fatal("Expected error to be ALREADY_EXISTS")
	}
}

func TestCreateItemNoExpiresAt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectExec("INSERT INTO item").WithArgs("id", "type", raw, nil).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewCreateItemCommand(service, &api.CreateItemRequest{
		Id:   "id",
		Type: "type",
		Data: data,
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.CreateItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestCreateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	inputTime := time.Now().UTC()
	mock.ExpectExec("INSERT INTO item").WithArgs("id", "type", raw, inputTime).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewCreateItemCommand(service, &api.CreateItemRequest{
		Id:        "id",
		Type:      "type",
		Data:      data,
		ExpiresAt: timestamppb.New(inputTime),
	})
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.CreateItemResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}
