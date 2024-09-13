package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/mysqlTestServer"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	server, err := mysqlTestServer.GetServer()
	if err != nil {
		log.Fatalf("could not run mysql test server: %v", err)
	}
	defer server.Close()
	db = server.Db
	schema, err := os.ReadFile("../../../migration/item.sql")
	if err != nil {
		log.Fatalf("could not read schema file: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("could not execute schema: %v", err)
	}

	m.Run()
}

func Test_CreateItem_ItemNoExpiry_ItemCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "1",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateItem_ItemWithExpiry_ItemCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "2",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateItem_ItemExists_ItemNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "3",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "3",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeDuplicateEntry {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_DeleteItem_ItemExists_ItemDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "4",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	result, err := q.DeleteItem(context.Background(), DeleteItemParams{
		ID:   "4",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not delete item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteItem_ItemDoesNotExist_ItemNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteItem(context.Background(), DeleteItemParams{
		ID:   "5",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not delete item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteItem_ItemExpired_ItemNotDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "6",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	result, err := q.DeleteItem(context.Background(), DeleteItemParams{
		ID:   "6",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not delete item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetItem_ItemExists_ItemReturned(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "7",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	item, err := q.GetItem(context.Background(), GetItemParams{
		ID:   "7",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not get item: %v", err)
	}
	if item.ID != "7" {
		t.Fatalf("expected id 7, got %s", item.ID)
	}
	if item.Type != "test" {
		t.Fatalf("expected type test, got %s", item.Type)
	}
	if item.ExpiresAt.Valid {
		t.Fatalf("expected expires_at to be null, got %v", item.ExpiresAt)
	}
}

func Test_GetItem_ItemDoesNotExist_ItemNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetItem(context.Background(), GetItemParams{
		ID:   "8",
		Type: "test",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetItem_ItemExpired_ItemNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "9",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.GetItem(context.Background(), GetItemParams{
		ID:   "9",
		Type: "test",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetItem_ItemNotExpired_ItemReturned(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "10",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	item, err := q.GetItem(context.Background(), GetItemParams{
		ID:   "10",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not get item: %v", err)
	}
	if item.ID != "10" {
		t.Fatalf("expected id 10, got %s", item.ID)
	}
	if item.Type != "test" {
		t.Fatalf("expected type test, got %s", item.Type)
	}
	if !item.ExpiresAt.Valid {
		t.Fatalf("expected expires_at to be valid, got %v", item.ExpiresAt)
	}
}

func Test_GetItems_NoType_GetAllItems(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "11",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "12",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	items, err := q.GetItems(context.Background(), GetItemsParams{
		Type:   sql.NullString{},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get items: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func Test_GetItems_WithType_GetItemsByType(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "13",
		Type:      "GetItems_WithType_GetItemsByType",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "14",
		Type:      "GetItems_WithType_GetItemsByType",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "15",
		Type:      "other",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	items, err := q.GetItems(context.Background(), GetItemsParams{
		Type: sql.NullString{
			Valid:  true,
			String: "GetItems_WithType_GetItemsByType",
		},
		Limit:  3,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get items: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].ID != "13" {
		t.Fatalf("expected id 13, got %s", items[0].ID)
	}
	if items[1].ID != "14" {
		t.Fatalf("expected id 14, got %s", items[1].ID)
	}
}

func Test_GetItems_WrongType_NoItems(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "16",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "17",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "18",
		Type:      "other",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	items, err := q.GetItems(context.Background(), GetItemsParams{
		Type: sql.NullString{
			Valid:  true,
			String: "wrong",
		},
		Limit:  3,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get items: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func Test_GetItems_ItemExpired_GetRemainingItems(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "19",
		Type:      "GetItems_ItemExpired_GetRemainingItems",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	_, err = q.CreateItem(context.Background(), CreateItemParams{
		ID:        "20",
		Type:      "GetItems_ItemExpired_GetRemainingItems",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	items, err := q.GetItems(context.Background(), GetItemsParams{
		Type: sql.NullString{
			Valid:  true,
			String: "GetItems_ItemExpired_GetRemainingItems",
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func Test_UpdateItem_ItemExists_ItemUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "21",
		Type:      "test",
		Data:      json.RawMessage(`{"key": "value"}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	result, err := q.UpdateItem(context.Background(), UpdateItemParams{
		ID:   "21",
		Type: "test",
		Data: json.RawMessage(`{"key": "new value"}`),
	})
	if err != nil {
		t.Fatalf("could not update item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_UpdateItem_ItemDoesNotExist_ItemNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateItem(context.Background(), UpdateItemParams{
		ID:   "22",
		Type: "test",
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateItem_ItemExpired_ItemNotUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateItem(context.Background(), CreateItemParams{
		ID:        "23",
		Type:      "test",
		Data:      json.RawMessage(`{"key": "value"}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create item: %v", err)
	}
	result, err := q.UpdateItem(context.Background(), UpdateItemParams{
		ID:   "23",
		Type: "test",
		Data: json.RawMessage(`{"key": "new value"}`),
	})
	if err != nil {
		t.Fatalf("could not update item: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}
