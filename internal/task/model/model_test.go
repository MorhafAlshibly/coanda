package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/mysqlTestServer"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var server *mysqlTestServer.Server

func TestMain(m *testing.M) {
	server = mysqlTestServer.NewServer("../../../migration/task.sql")
	defer server.Close()
	m.Run()
}

func Test_CreateTask_TaskNoExpiry_TaskCreated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "1",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateTask_TaskWithExpiry_TaskCreated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "2",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateTask_TaskExists_TaskNotCreated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "3",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
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

func Test_DeleteTask_TaskExists_TaskDeleted(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "4",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	result, err := q.DeleteTask(context.Background(), DeleteTaskParams{
		ID:   "4",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not delete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteTask_TaskDoesNotExist_TaskNotDeleted(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.DeleteTask(context.Background(), DeleteTaskParams{
		ID:   "5",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not delete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteTask_TaskExpired_TaskNotDeleted(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "6",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	result, err := q.DeleteTask(context.Background(), DeleteTaskParams{
		ID:   "6",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not delete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetTask_TaskExists_TaskReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "7",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	task, err := q.GetTask(context.Background(), GetTaskParams{
		ID:   "7",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not get task: %v", err)
	}
	if task.ID != "7" {
		t.Fatalf("expected id 7, got %s", task.ID)
	}
	if task.Type != "test" {
		t.Fatalf("expected type test, got %s", task.Type)
	}
	if task.ExpiresAt.Valid {
		t.Fatalf("expected expires_at to be null, got %v", task.ExpiresAt)
	}
}

func Test_GetTask_TaskDoesNotExist_TaskNotReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.GetTask(context.Background(), GetTaskParams{
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

func Test_GetTask_TaskExpired_TaskNotReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "9",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.GetTask(context.Background(), GetTaskParams{
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

func Test_GetTask_TaskNotExpired_TaskReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "10",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	task, err := q.GetTask(context.Background(), GetTaskParams{
		ID:   "10",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not get task: %v", err)
	}
	if task.ID != "10" {
		t.Fatalf("expected id 10, got %s", task.ID)
	}
	if task.Type != "test" {
		t.Fatalf("expected type test, got %s", task.Type)
	}
	if !task.ExpiresAt.Valid {
		t.Fatalf("expected expires_at to be valid, got %v", task.ExpiresAt)
	}
}

func Test_GetTasks_NoType_GetAllTasks(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "11",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "12",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	tasks, err := q.GetTasks(context.Background(), GetTasksParams{
		Type:   sql.NullString{},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func Test_GetTasks_WithType_GetTasksByType(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "13",
		Type:      "GetTasks_WithType_GetTasksByType",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "14",
		Type:      "GetTasks_WithType_GetTasksByType",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "15",
		Type:      "other",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	tasks, err := q.GetTasks(context.Background(), GetTasksParams{
		Type: sql.NullString{
			Valid:  true,
			String: "GetTasks_WithType_GetTasksByType",
		},
		Limit:  3,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != "13" {
		t.Fatalf("expected id 13, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "14" {
		t.Fatalf("expected id 14, got %s", tasks[1].ID)
	}
}

func Test_GetTasks_WrongType_NoTasks(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "16",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "17",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "18",
		Type:      "other",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	tasks, err := q.GetTasks(context.Background(), GetTasksParams{
		Type: sql.NullString{
			Valid:  true,
			String: "wrong",
		},
		Limit:  3,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get tasks: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func Test_GetTasks_TaskExpired_GetRemainingTasks(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "19",
		Type:      "GetTasks_TaskExpired_GetRemainingTasks",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "20",
		Type:      "GetTasks_TaskExpired_GetRemainingTasks",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	tasks, err := q.GetTasks(context.Background(), GetTasksParams{
		Type: sql.NullString{
			Valid:  true,
			String: "GetTasks_TaskExpired_GetRemainingTasks",
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func Test_UpdateTask_TaskExists_TaskUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "21",
		Type:      "test",
		Data:      json.RawMessage(`{"key": "value"}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	result, err := q.UpdateTask(context.Background(), UpdateTaskParams{
		ID:   "21",
		Type: "test",
		Data: json.RawMessage(`{"key": "new value"}`),
	})
	if err != nil {
		t.Fatalf("could not update task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_UpdateTask_TaskDoesNotExist_TaskNotUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.UpdateTask(context.Background(), UpdateTaskParams{
		ID:   "22",
		Type: "test",
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateTask_TaskExpired_TaskNotUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "23",
		Type:      "test",
		Data:      json.RawMessage(`{"key": "value"}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	result, err := q.UpdateTask(context.Background(), UpdateTaskParams{
		ID:   "23",
		Type: "test",
		Data: json.RawMessage(`{"key": "new value"}`),
	})
	if err != nil {
		t.Fatalf("could not update task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_CompleteTask_TaskExists_TaskCompleted(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "24",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	result, err := q.CompleteTask(context.Background(), CompleteTaskParams{
		ID:   "24",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not complete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CompleteTask_TaskDoesNotExist_TaskNotCompleted(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.CompleteTask(context.Background(), CompleteTaskParams{
		ID:   "25",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not complete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_CompleteTask_TaskExpired_TaskNotCompleted(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "26",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	result, err := q.CompleteTask(context.Background(), CompleteTaskParams{
		ID:   "26",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not complete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_CompleteTask_TaskCompleted_TaskNotCompletedAgain(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.CreateTask(context.Background(), CreateTaskParams{
		ID:        "27",
		Type:      "test",
		Data:      json.RawMessage(`{}`),
		ExpiresAt: sql.NullTime{},
	})
	if err != nil {
		t.Fatalf("could not create task: %v", err)
	}
	_, err = q.CompleteTask(context.Background(), CompleteTaskParams{
		ID:   "27",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not complete task: %v", err)
	}
	task, err := q.GetTask(context.Background(), GetTaskParams{
		ID:   "27",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not get task: %v", err)
	}
	completedAt := task.CompletedAt
	result, err := q.CompleteTask(context.Background(), CompleteTaskParams{
		ID:   "27",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not complete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
	task, err = q.GetTask(context.Background(), GetTaskParams{
		ID:   "27",
		Type: "test",
	})
	if err != nil {
		t.Fatalf("could not get task: %v", err)
	}
	if task.CompletedAt != completedAt {
		t.Fatalf("expected completed_at to be %v, got %v", completedAt, task.CompletedAt)
	}
}
