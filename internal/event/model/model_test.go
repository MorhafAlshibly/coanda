package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"reflect"
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
	schema, err := os.ReadFile("../../../migration/event.sql")
	if err != nil {
		log.Fatalf("could not read schema file: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("could not execute schema: %v", err)
	}

	m.Run()
}

func Test_CreateEvent_Event_EventCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateEvent_EventExists_EventNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event1",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	_, err = q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event1",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
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

func Test_CreateEventRound_EventRound_EventRoundCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event2",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateEventRound_EventRoundExists_EventRoundNotCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event3",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round1",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round1",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now(),
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

func Test_CreateEventRound_EventRoundEndsAtSameTime_EventRoundNotCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event4",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	endedAt := time.Now()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round2",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: endedAt,
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round3",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: endedAt,
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

func Test_CreateEventRound_EventDoesNotExist_EventRoundNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: 999999,
		Name:    "round4",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now(),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_CreateOrUpdateEventUser_EventUser_EventUserCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event5",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  1,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateOrUpdateEventUser_EventUserExists_EventUserUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event6",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  2,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result2, err := q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  2,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	rowsAffected, err := result2.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 2 {
		t.Fatalf("expected 2 row affected, got %d", rowsAffected)
	}
	eventUserId2, err := result2.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	if eventUserId != eventUserId2 {
		t.Fatalf("expected event user id to be the same, got %d and %d", eventUserId, eventUserId2)
	}
}

func Test_CreateOrUpdateEventUser_EventDoesNotExist_EventUserNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: 99999,
		UserID:  3,
		Data:    json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_GetEventRoundUserByEventUserId_EventRoundUser_EventRoundUserReturned(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event7",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round5",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  4,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId),
		Result:       0,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	row, err := q.GetEventRoundUserByEventUserId(context.Background(), uint64(eventUserId))
	if err != nil {
		t.Fatalf("could not get event round user: %v", err)
	}
	if row.EventUserID != uint64(eventUserId) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId, row.EventUserID)
	}
	if row.EventRoundID != uint64(eventRoundId) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId, row.EventRoundID)
	}
	if row.Result != 0 {
		t.Fatalf("expected result to be 0, got %d", row.Result)
	}
}

func Test_GetEventRoundUserByEventUserId_EventRoundUserDoesNotExist_EventRoundUserNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEventRoundUserByEventUserId(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func Test_GetEventRoundUserByEventUserId_EventRoundUserInAnotherRound_EventRoundUserReturned(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event8",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round6",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round7",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now().Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  5,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId1),
		Result:       0,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUser, err := q.GetEventRoundUserByEventUserId(context.Background(), uint64(eventUserId))
	if err != nil {
		t.Fatalf("could not get event round user: %v", err)
	}
	if eventRoundUser.EventRoundID != uint64(eventRoundId1) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId1, eventRoundUser.EventRoundID)
	}
}

func Test_CreateEventRoundUser_EventRoundUser_EventRoundUserCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event9",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  6,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round8",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId),
		Result:       0,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateEventRoundUser_EventRoundUserExists_EventRoundUserNotCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event10",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  7,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round9",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	_, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId),
		Result:       0,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	_, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
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

func Test_CreateEventRoundUser_EventRoundDoesNotExist_EventRoundUserNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  9999,
		EventRoundID: 9999,
		Result:       0,
		Data:         json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_UpdateEventRoundUserResult_UpdateResultAndData_EventRoundUserResultUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event12",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(id),
		UserID:  9,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(id),
		Name:    "round11",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{}`),
		EndedAt: time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId),
		Result:       0,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	result, err = q.UpdateEventRoundUserResult(context.Background(), UpdateEventRoundUserResultParams{
		EventUserID:  uint64(eventUserId),
		EventRoundID: uint64(eventRoundId),
		Result:       1,
		Data:         json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event round user result: %v", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	row, err := q.GetEventRoundUserByEventUserId(context.Background(), uint64(eventUserId))
	if err != nil {
		t.Fatalf("could not get event round user: %v", err)
	}
	if row.Result != 1 {
		t.Fatalf("expected result to be 1, got %d", row.Result)
	}
	if string(row.Data) != `{"key": "value"}` {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
}

func Test_UpdateEventRoundUserResult_EventRoundUserDoesNotExist_EventRoundUserResultNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateEventRoundUserResult(context.Background(), UpdateEventRoundUserResultParams{
		EventUserID: 999999,
		Result:      0,
		Data:        json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not update event round user result: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 row affected, got %d", rowsAffected)
	}
}

func Test_GetEvent_ById_EventReturned(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event14",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	row, err := q.GetEvent(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event: %v", err)
	}
	if row.ID != uint64(id) {
		t.Fatalf("expected id to be %d, got %d", id, row.ID)
	}
	if row.Name != "event14" {
		t.Fatalf("expected name to be event14, got %s", row.Name)
	}
}

func Test_GetEvent_ByName_EventReturned(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event15",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	_, err = result.LastInsertId()
	row, err := q.GetEvent(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event15", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event: %v", err)
	}
	if row.Name != "event15" {
		t.Fatalf("expected name to be event15, got %s", row.Name)
	}
}

func Test_GetEvent_ByIdEventDoesNotExist_EventNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEvent(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: 999999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetEvent_ByNameEventDoesNotExist_EventNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEvent(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event16", Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEvent_ById_EventDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event17",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	_, err = q.DeleteEvent(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event: %v", err)
	}
	_, err = q.GetEvent(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEvent_ByName_EventDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event18",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	_, err = q.DeleteEvent(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event18", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event: %v", err)
	}
	_, err = q.GetEvent(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event18", Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEvent_ByIdEventDoesNotExist_EventNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteEvent(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteEvent_ByNameEventDoesNotExist_EventNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteEvent(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event19", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetEventWithRound_ByID_EventWithRoundReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event20",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round13",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	eventWithRounds, err := q.GetEventWithRound(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: eventId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event with round: %v", err)
	}
	if len(eventWithRounds) != 1 {
		t.Fatalf("expected 1 event and 1 round, got %d", len(eventWithRounds))
	}
	if eventWithRounds[0].ID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventWithRounds[0].ID)
	}
	if eventWithRounds[0].Name != "event20" {
		t.Fatalf("expected event name to be event20, got %s", eventWithRounds[0].Name)
	}
	if eventWithRounds[0].CurrentRoundID == nil {
		t.Fatalf("expected event current round id to be not nil, got nil")
	}
	if *eventWithRounds[0].CurrentRoundID != uint64(eventRoundId) {
		t.Fatalf("expected event current round id to be %d, got %d", eventRoundId, *eventWithRounds[0].CurrentRoundID)
	}
	if eventWithRounds[0].CurrentRoundName == nil {
		t.Fatalf("expected event current round name to be not nil, got nil")
	}
	if *eventWithRounds[0].CurrentRoundName != "round13" {
		t.Fatalf("expected event current round name to be round13, got %s", *eventWithRounds[0].CurrentRoundName)
	}
	if !reflect.DeepEqual(eventWithRounds[0].Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected event data to be {}, got %s", eventWithRounds[0].Data)
	}
	if eventWithRounds[0].RoundID.Int64 != eventRoundId {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId, eventWithRounds[0].RoundID.Int64)
	}
	if eventWithRounds[0].RoundName.String != "round13" {
		t.Fatalf("expected event round name to be round13, got %s", eventWithRounds[0].RoundName.String)
	}
	if !reflect.DeepEqual(eventWithRounds[0].RoundScoring, json.RawMessage(`{"scoring": [1, 2, 4]}`)) {
		t.Fatalf("expected event round scoring to be {\"scoring\": [1, 2, 4]}, got %s", eventWithRounds[0].RoundScoring)
	}
	if !reflect.DeepEqual(eventWithRounds[0].RoundData, json.RawMessage(`{}`)) {
		t.Fatalf("expected event round data to be {}, got %s", eventWithRounds[0].RoundData)
	}
	if !eventWithRounds[0].RoundEndedAt.Time.Equal(startedAt.Add(1 * time.Hour)) {
		t.Fatalf("expected event round ended at to be %v, got %v", startedAt.Add(1*time.Hour), eventWithRounds[0].RoundEndedAt)
	}
	if !eventWithRounds[0].StartedAt.Equal(startedAt) {
		t.Fatalf("expected event started at to be %v, got %v", startedAt, eventWithRounds[0].StartedAt)
	}
}

func Test_GetEventWithRound_ByName_EventWithRoundReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event21",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round14",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	eventWithRounds, err := q.GetEventWithRound(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event21", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event with round: %v", err)
	}
	if len(eventWithRounds) != 1 {
		t.Fatalf("expected 1 event and 1 round, got %d", len(eventWithRounds))
	}
	if eventWithRounds[0].ID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventWithRounds[0].ID)
	}
	if eventWithRounds[0].Name != "event21" {
		t.Fatalf("expected event name to be event21, got %s", eventWithRounds[0].Name)
	}
	if *eventWithRounds[0].CurrentRoundID != uint64(eventRoundId) {
		t.Fatalf("expected event current round id to be %d, got %d", eventRoundId, *eventWithRounds[0].CurrentRoundID)
	}
}

func Test_GetEventWithRound_ByIDEventDoesNotExist_EventNotReturned(t *testing.T) {
	q := New(db)
	eventRounds, err := q.GetEventWithRound(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event with round: %v", err)
	}
	if len(eventRounds) != 0 {
		t.Fatalf("expected 0 event and 0 round, got %d", len(eventRounds))
	}
}

func Test_GetEventWithRound_ByNameEventDoesNotExist_EventNotReturned(t *testing.T) {
	q := New(db)
	eventWithRounds, err := q.GetEventWithRound(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event22", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event with round: %v", err)
	}
	if len(eventWithRounds) != 0 {
		t.Fatalf("expected 0 event and 0 round, got %d", len(eventWithRounds))
	}
}

func Test_GetEventWithRound_ByIDMultipleEventRounds_EventReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event23",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round15",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round16",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventWithRounds, err := q.GetEventWithRound(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: eventId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event with round: %v", err)
	}
	if len(eventWithRounds) != 2 {
		t.Fatalf("expected 2 rounds, got %d", len(eventWithRounds))
	}
	if eventWithRounds[0].ID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventWithRounds[0].ID)
	}
	if eventWithRounds[0].Name != "event23" {
		t.Fatalf("expected event name to be event23, got %s", eventWithRounds[0].Name)
	}
	if *eventWithRounds[0].CurrentRoundID != uint64(eventRoundId1) {
		t.Fatalf("expected event current round id to be %d, got %d", eventRoundId1, *eventWithRounds[0].CurrentRoundID)
	}
}

func Test_GetEventWithRound_ByIDEventEnded_EventReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event24",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt.Add(-2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round17",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(-1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventWithRounds, err := q.GetEventWithRound(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: eventId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event with round: %v", err)
	}
	if len(eventWithRounds) != 1 {
		t.Fatalf("expected 1 event and 1 round, got %d", len(eventWithRounds))
	}
	if eventWithRounds[0].ID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventWithRounds[0].ID)
	}
	if eventWithRounds[0].Name != "event24" {
		t.Fatalf("expected event name to be event24, got %s", eventWithRounds[0].Name)
	}
	if eventWithRounds[0].CurrentRoundID != nil {
		t.Fatalf("expected event current round id to be nil, got %d", *eventWithRounds[0].CurrentRoundID)
	}
}

func Test_GetEventLeaderboard_ByID_EventLeaderboardReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event25",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  11,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId1, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  12,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId2, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round18",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId1),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId2),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	leaderboard, err := q.GetEventLeaderboard(context.Background(), GetEventLeaderboardParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event leaderboard: %v", err)
	}
	if len(leaderboard) != 2 {
		t.Fatalf("expected 2 leaderboard entries, got %d", len(leaderboard))
	}
	if leaderboard[0].UserID != 11 {
		t.Fatalf("expected user id to be 11, got %d", leaderboard[0].UserID)
	}
	if leaderboard[0].Score != 4 {
		t.Fatalf("expected score to be 4, got %d", leaderboard[0].Score)
	}
	if leaderboard[0].Ranking != 1 {
		t.Fatalf("expected ranking to be 1, got %d", leaderboard[0].Ranking)
	}
	if leaderboard[1].UserID != 12 {
		t.Fatalf("expected user id to be 12, got %d", leaderboard[1].UserID)
	}
	if leaderboard[1].Score != 1 {
		t.Fatalf("expected score to be 1, got %d", leaderboard[1].Score)
	}
	if leaderboard[1].Ranking != 2 {
		t.Fatalf("expected ranking to be 2, got %d", leaderboard[1].Ranking)
	}
}

func Test_GetEventLeaderboard_ByName_EventLeaderboardReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event26",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  13,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId1, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  14,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId2, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round19",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId1),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId2),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	leaderboard, err := q.GetEventLeaderboard(context.Background(), GetEventLeaderboardParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event26", Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event leaderboard: %v", err)
	}
	if len(leaderboard) != 2 {
		t.Fatalf("expected 2 leaderboard entries, got %d", len(leaderboard))
	}
	if leaderboard[0].UserID != 13 {
		t.Fatalf("expected user id to be 13, got %d", leaderboard[0].UserID)
	}
	if leaderboard[0].Score != 4 {
		t.Fatalf("expected score to be 4, got %d", leaderboard[0].Score)
	}
	if leaderboard[0].Ranking != 1 {
		t.Fatalf("expected ranking to be 1, got %d", leaderboard[0].Ranking)
	}
	if leaderboard[1].UserID != 14 {
		t.Fatalf("expected user id to be 14, got %d", leaderboard[1].UserID)
	}
	if leaderboard[1].Score != 1 {
		t.Fatalf("expected score to be 1, got %d", leaderboard[1].Score)
	}
	if leaderboard[1].Ranking != 2 {
		t.Fatalf("expected ranking to be 2, got %d", leaderboard[1].Ranking)
	}
}

func Test_GetEventLeaderboard_ByIDEventDoesNotExist_EventLeaderboardNotReturned(t *testing.T) {
	q := New(db)
	eventLeaderboard, err := q.GetEventLeaderboard(context.Background(), GetEventLeaderboardParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event leaderboard: %v", err)
	}
	if len(eventLeaderboard) != 0 {
		t.Fatalf("expected 0 leaderboard entries, got %d", len(eventLeaderboard))
	}
}

func Test_GetEventLeaderboard_ByNameEventDoesNotExist_EventLeaderboardNotReturned(t *testing.T) {
	q := New(db)
	eventLeaderboard, err := q.GetEventLeaderboard(context.Background(), GetEventLeaderboardParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event27", Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event leaderboard: %v", err)
	}
	if len(eventLeaderboard) != 0 {
		t.Fatalf("expected 0 leaderboard entries, got %d", len(eventLeaderboard))
	}
}

func Test_GetEventLeaderboard_ByIDMultipleEventRounds_EventLeaderboardReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event28",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  15,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId1, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  16,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId2, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round20",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventUserId1),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventUserId2),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round21",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [10,9]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId2, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId2),
		EventUserID:  uint64(eventUserId1),
		Result:       8,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId2),
		EventUserID:  uint64(eventUserId2),
		Result:       9,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	leaderboard, err := q.GetEventLeaderboard(context.Background(), GetEventLeaderboardParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event leaderboard: %v", err)
	}
	if len(leaderboard) != 2 {
		t.Fatalf("expected 2 leaderboard entries, got %d", len(leaderboard))
	}
	if leaderboard[0].UserID != 15 {
		t.Fatalf("expected user id to be 15, got %d", leaderboard[0].UserID)
	}
	if leaderboard[0].Score != 14 {
		t.Fatalf("expected score to be 14, got %d", leaderboard[0].Score)
	}
	if leaderboard[0].Ranking != 1 {
		t.Fatalf("expected ranking to be 1, got %d", leaderboard[0].Ranking)
	}
	if leaderboard[1].UserID != 16 {
		t.Fatalf("expected user id to be 16, got %d", leaderboard[1].UserID)
	}
	if leaderboard[1].Score != 10 {
		t.Fatalf("expected score to be 10, got %d", leaderboard[1].Score)
	}
	if leaderboard[1].Ranking != 2 {
		t.Fatalf("expected ranking to be 2, got %d", leaderboard[1].Ranking)
	}
}

func Test_UpdateEvent_ByID_EventUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event29",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	id, err := result.LastInsertId()
	_, err = q.UpdateEvent(context.Background(), UpdateEventParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: id, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event: %v", err)
	}
	row, err := q.GetEvent(context.Background(), GetEventParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
}

func Test_UpdateEvent_ByName_EventUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event30",
		Data:      json.RawMessage(`{}`),
		StartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	_, err = q.UpdateEvent(context.Background(), UpdateEventParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event30", Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event: %v", err)
	}
	row, err := q.GetEvent(context.Background(), GetEventParams{
		Name: sql.NullString{String: "event30", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
}

func Test_UpdateEvent_ByIDEventDoesNotExist_EventNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateEvent(context.Background(), UpdateEventParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateEvent_ByNameEventDoesNotExist_EventNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateEvent(context.Background(), UpdateEventParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event31", Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetEventRound_ByEventID_CurrentEventRoundReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event32",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round22",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round23",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRound, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if eventRound.ID != uint64(eventRoundId1) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId1, eventRound.ID)
	}
	if eventRound.Name != "round22" {
		t.Fatalf("expected event round name to be round22, got %s", eventRound.Name)
	}
	if !reflect.DeepEqual(eventRound.Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected event round data to be {}, got %s", eventRound.Data)
	}
	if !reflect.DeepEqual(eventRound.Scoring, json.RawMessage(`{"scoring": [1, 2, 4]}`)) {
		t.Fatalf("expected event round scoring to be {\"scoring\": [1, 2, 4]}, got %s", eventRound.Scoring)
	}
	if !eventRound.EndedAt.Equal(startedAt.Add(1 * time.Hour)) {
		t.Fatalf("expected event round ended at to be %v, got %v", startedAt.Add(1*time.Hour), eventRound.EndedAt)
	}
}

func Test_GetEventRound_ByEventName_CurrentEventRoundReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event33",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round24",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round25",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRound, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event33", Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if eventRound.Name != "round24" {
		t.Fatalf("expected event round name to be round24, got %s", eventRound.Name)
	}
	if !reflect.DeepEqual(eventRound.Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected event round data to be {}, got %s", eventRound.Data)
	}
	if !reflect.DeepEqual(eventRound.Scoring, json.RawMessage(`{"scoring": [1, 2, 4]}`)) {
		t.Fatalf("expected event round scoring to be {\"scoring\": [1, 2, 4]}, got %s", eventRound.Scoring)
	}
	if !eventRound.EndedAt.Equal(startedAt.Add(1 * time.Hour)) {
		t.Fatalf("expected event round ended at to be %v, got %v", startedAt.Add(1*time.Hour), eventRound.EndedAt)
	}
}

func Test_GetEventRound_ByEventIDEventDoesNotExist_EventRoundNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetEventRound_ByEventNameEventDoesNotExist_EventRoundNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event34", Valid: true},
		},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetEventRound_ByRoundID_RoundReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event35",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round26",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	eventRound, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		ID: sql.NullInt64{Int64: eventRoundId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if eventRound.Name != "round26" {
		t.Fatalf("expected event round name to be round26, got %s", eventRound.Name)
	}
	if !reflect.DeepEqual(eventRound.Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected event round data to be {}, got %s", eventRound.Data)
	}
	if !reflect.DeepEqual(eventRound.Scoring, json.RawMessage(`{"scoring": [1, 2, 4]}`)) {
		t.Fatalf("expected event round scoring to be {\"scoring\": [1, 2, 4]}, got %s", eventRound.Scoring)
	}
	if !eventRound.EndedAt.Equal(startedAt.Add(1 * time.Hour)) {
		t.Fatalf("expected event round ended at to be %v, got %v", startedAt.Add(1*time.Hour), eventRound.EndedAt)
	}
}

func Test_GetEventRound_ByRoundName_RoundReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event36",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round27",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRound, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event36", Valid: true},
		},
		Name: sql.NullString{String: "round27", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if eventRound.Name != "round27" {
		t.Fatalf("expected event round name to be round27, got %s", eventRound.Name)
	}
	if !reflect.DeepEqual(eventRound.Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected event round data to be {}, got %s", eventRound.Data)
	}
	if !reflect.DeepEqual(eventRound.Scoring, json.RawMessage(`{"scoring": [1, 2, 4]}`)) {
		t.Fatalf("expected event round scoring to be {\"scoring\": [1, 2, 4]}, got %s", eventRound.Scoring)
	}
	if !eventRound.EndedAt.Equal(startedAt.Add(1 * time.Hour)) {
		t.Fatalf("expected event round ended at to be %v, got %v", startedAt.Add(1*time.Hour), eventRound.EndedAt)
	}
}

func Test_GetEventRound_ByRoundIDRoundDoesNotExist_RoundNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		ID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetEventRoundLeaderboard_ByEventID_RoundLeaderboardReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event37",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  17,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId1, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  18,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId2, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round28",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round29",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId1),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId1, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId2),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	leaderboard, err := q.GetEventRoundLeaderboard(context.Background(), GetEventRoundLeaderboardParams{
		EventRound: GetEventRoundParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: eventId, Valid: true},
			},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round leaderboard: %v", err)
	}
	if len(leaderboard) != 2 {
		t.Fatalf("expected 2 leaderboard entries, got %d", len(leaderboard))
	}
	if leaderboard[0].ID != uint64(eventRoundUserId1) {
		t.Fatalf("expected round user id to be %d, got %d", eventRoundUserId1, leaderboard[0].ID)
	}
	if leaderboard[0].EventID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, leaderboard[0].EventID)
	}
	if leaderboard[0].RoundName != "round28" {
		t.Fatalf("expected round name to be round28, got %s", leaderboard[0].RoundName)
	}
	if leaderboard[0].EventUserID != uint64(eventUserId2) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId1, leaderboard[0].EventUserID)
	}
	if leaderboard[0].EventRoundID != uint64(eventRoundId) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId, leaderboard[0].EventRoundID)
	}
	if leaderboard[0].Result != 2 {
		t.Fatalf("expected result to be 2, got %d", leaderboard[0].Result)
	}
	if leaderboard[1].Score != 4 {
		t.Fatalf("expected score to be 4, got %d", leaderboard[1].Score)
	}
	if leaderboard[1].Ranking != 1 {
		t.Fatalf("expected ranking to be 1, got %d", leaderboard[1].Ranking)
	}
}

func Test_GetEventRoundLeaderboard_ByEventName_RoundLeaderboardReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event38",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  19,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId1, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  20,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId2, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round30",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round31",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	//eventRoundId2, err := result.LastInsertId()
	_, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventUserId1),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventUserId2),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	leaderboard, err := q.GetEventRoundLeaderboard(context.Background(), GetEventRoundLeaderboardParams{
		EventRound: GetEventRoundParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event38", Valid: true},
			},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round leaderboard: %v", err)
	}
	if len(leaderboard) != 2 {
		t.Fatalf("expected 2 leaderboard entries, got %d", len(leaderboard))
	}
}

func Test_GetEventRoundLeaderboard_ByRoundIDRoundDoesNotExist_RoundLeaderboardNotReturned(t *testing.T) {
	q := New(db)
	eventRoundLeaderboard, err := q.GetEventRoundLeaderboard(context.Background(), GetEventRoundLeaderboardParams{
		EventRound: GetEventRoundParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round leaderboard: %v", err)
	}
	if len(eventRoundLeaderboard) != 0 {
		t.Fatalf("expected 0 leaderboard entries, got %d", len(eventRoundLeaderboard))
	}
}

func Test_GetEventRoundLeaderboard_ByRoundNameRoundDoesNotExist_RoundLeaderboardNotReturned(t *testing.T) {
	q := New(db)
	eventRoundLeaderboard, err := q.GetEventRoundLeaderboard(context.Background(), GetEventRoundLeaderboardParams{
		EventRound: GetEventRoundParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event39", Valid: true},
			},
			Name: sql.NullString{String: "round32", Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round leaderboard: %v", err)
	}
	if len(eventRoundLeaderboard) != 0 {
		t.Fatalf("expected 0 leaderboard entries, got %d", len(eventRoundLeaderboard))
	}
}

func Test_GetEventRoundLeaderboard_ByEventIDAllRoundsEnded_NoLeaderboardReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event40",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt.Add(-3 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  21,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  22,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round33",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(-1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round34",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(-2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId2, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventRoundId1),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventRoundId2),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	_, err = q.GetEventRoundLeaderboard(context.Background(), GetEventRoundLeaderboardParams{
		EventRound: GetEventRoundParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: eventId, Valid: true},
			},
		},
		Limit:  2,
		Offset: 0,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_UpdateEventRound_ByEventID_RoundUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event41",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round35",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	_, err = q.UpdateEventRound(context.Background(), UpdateEventRoundParams{
		EventRound: GetEventRoundParams{
			ID: sql.NullInt64{Int64: eventRoundId, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event round: %v", err)
	}
	row, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
}

func Test_UpdateEventRound_ByEventName_RoundUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event42",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round36",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round37",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.UpdateEventRound(context.Background(), UpdateEventRoundParams{
		EventRound: GetEventRoundParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event42", Valid: true},
			},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event round: %v", err)
	}
	row, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event42", Valid: true},
		},
		Name: sql.NullString{String: "round36", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
	row2, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event42", Valid: true},
		},
		Name: sql.NullString{String: "round37", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if !reflect.DeepEqual(row2.Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected data to be {}, got %s", row2.Data)
	}
}

func Test_UpdateEventRound_ByRoundID_RoundUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event43",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round38",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	_, err = q.UpdateEventRound(context.Background(), UpdateEventRoundParams{
		EventRound: GetEventRoundParams{
			ID: sql.NullInt64{Int64: eventRoundId, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event round: %v", err)
	}
	row, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		ID: sql.NullInt64{Int64: eventRoundId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
}

func Test_UpdateEventRound_ByRoundName_RoundUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event44",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round39",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round40",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [1, 2, 4]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	_, err = q.UpdateEventRound(context.Background(), UpdateEventRoundParams{
		EventRound: GetEventRoundParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event44", Valid: true},
			},
			Name: sql.NullString{String: "round39", Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update event round: %v", err)
	}
	row, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event44", Valid: true},
		},
		Name: sql.NullString{String: "round39", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", row.Data)
	}
	row2, err := q.GetEventRound(context.Background(), GetEventRoundParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event44", Valid: true},
		},
		Name: sql.NullString{String: "round40", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event round: %v", err)
	}
	if !reflect.DeepEqual(row2.Data, json.RawMessage(`{}`)) {
		t.Fatalf("expected data to be {}, got %s", row2.Data)
	}
}

func Test_GetEventUser_ByID_EventUserReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event45",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  23,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	eventUser, err := q.GetEventUser(context.Background(), GetEventUserParams{
		ID: sql.NullInt64{Int64: eventUserId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event user: %v", err)
	}
	if eventUser.EventID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventUser.EventID)
	}
	if eventUser.UserID != 23 {
		t.Fatalf("expected user id to be 23, got %d", eventUser.UserID)
	}
	if !reflect.DeepEqual(eventUser.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", eventUser.Data)
	}
}

func Test_GetEventUser_ByEventIDAndUserID_EventUserReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event46",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  24,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUser, err := q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: 1, Valid: true},
		},
		UserID: sql.NullInt64{Int64: 24, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event user: %v", err)
	}
	if eventUser.EventID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventUser.EventID)
	}
	if eventUser.UserID != 24 {
		t.Fatalf("expected user id to be 24, got %d", eventUser.UserID)
	}
	if !reflect.DeepEqual(eventUser.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", eventUser.Data)
	}
}

func Test_GetEventUser_ByEventNameAndUserID_EventUserReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event47",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  25,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	eventUser, err := q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event47", Valid: true},
		},
		UserID: sql.NullInt64{Int64: 25, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event user: %v", err)
	}
	if eventUser.EventID != uint64(eventId) {
		t.Fatalf("expected event id to be %d, got %d", eventId, eventUser.EventID)
	}
	if eventUser.ID != uint64(eventUserId) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId, eventUser.ID)
	}
	if eventUser.UserID != 25 {
		t.Fatalf("expected user id to be 25, got %d", eventUser.UserID)
	}
	if !reflect.DeepEqual(eventUser.Data, json.RawMessage(`{"key": "value"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value\"}, got %s", eventUser.Data)
	}
}

func Test_GetEventUser_ByEventIDAndUserIDEventDoesNotExist_EventUserNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		UserID: sql.NullInt64{Int64: 26, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetEventUser_ByEventNameAndUserIDEventDoesNotExist_EventUserNotReturned(t *testing.T) {
	q := New(db)
	_, err := q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event48", Valid: true},
		},
		UserID: sql.NullInt64{Int64: 27, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_GetEventUser_ByEventIDAndUserIDUserDoesNotExist_EventUserNotReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event49",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  28,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	_, err = q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
		UserID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetEventRoundUsers_ByEventUserID_EventRoundUsersReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event50",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  29,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round41",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round42",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [2,3]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId2, err := result.LastInsertId()
	_, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId2),
		EventUserID:  uint64(eventUserId),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUsers, err := q.GetEventRoundUsers(context.Background(), GetEventRoundUsersParams{
		EventUser: GetEventUserParams{
			ID: sql.NullInt64{Int64: eventUserId, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round users: %v", err)
	}
	if len(eventRoundUsers) != 2 {
		t.Fatalf("expected 2 event round users, got %d", len(eventRoundUsers))
	}
	if eventRoundUsers[0].EventUserID != uint64(eventUserId) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId, eventRoundUsers[0].EventUserID)
	}
	if eventRoundUsers[1].EventUserID != uint64(eventUserId) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId, eventRoundUsers[1].EventUserID)
	}
	if eventRoundUsers[0].EventRoundID != uint64(eventRoundId2) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId2, eventRoundUsers[0].EventRoundID)
	}
	if eventRoundUsers[1].EventRoundID != uint64(eventRoundId1) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId1, eventRoundUsers[1].EventRoundID)
	}
	if eventRoundUsers[0].Result != 2 {
		t.Fatalf("expected result to be 2, got %d", eventRoundUsers[0].Result)
	}
	if eventRoundUsers[1].Result != 1 {
		t.Fatalf("expected result to be 1, got %d", eventRoundUsers[1].Result)
	}
	if eventRoundUsers[0].Ranking != 1 {
		t.Fatalf("expected ranking to be 1, got %d", eventRoundUsers[0].Ranking)
	}
	if eventRoundUsers[1].Ranking != 2 {
		t.Fatalf("expected ranking to be 2, got %d", eventRoundUsers[1].Ranking)
	}
}

func Test_GetEventRoundUsers_ByEventUserIDEventUserDoesNotExist_EventRoundUsersNotReturned(t *testing.T) {
	q := New(db)
	eventRoundUsers, err := q.GetEventRoundUsers(context.Background(), GetEventRoundUsersParams{
		EventUser: GetEventUserParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round users: %v", err)
	}
	if len(eventRoundUsers) != 0 {
		t.Fatalf("expected 0 event round users, got %d", len(eventRoundUsers))
	}
}

func Test_GetEventRoundUsers_ByEventIDAndUserID_EventRoundUsersReturned(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event51",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  30,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round43",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId1, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round44",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [2,3]}`),
		EndedAt: startedAt.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId2, err := result.LastInsertId()
	_, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId1),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId2),
		EventUserID:  uint64(eventUserId),
		Result:       2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUsers, err := q.GetEventRoundUsers(context.Background(), GetEventRoundUsersParams{
		EventUser: GetEventUserParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: eventId, Valid: true},
			},
			UserID: sql.NullInt64{Int64: 30, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get event round users: %v", err)
	}
	if len(eventRoundUsers) != 2 {
		t.Fatalf("expected 2 event round users, got %d", len(eventRoundUsers))
	}
	if eventRoundUsers[0].EventUserID != uint64(eventUserId) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId, eventRoundUsers[0].EventUserID)
	}
	if eventRoundUsers[1].EventUserID != uint64(eventUserId) {
		t.Fatalf("expected event user id to be %d, got %d", eventUserId, eventRoundUsers[1].EventUserID)
	}
	if eventRoundUsers[0].EventRoundID != uint64(eventRoundId2) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId2, eventRoundUsers[0].EventRoundID)
	}
	if eventRoundUsers[1].EventRoundID != uint64(eventRoundId1) {
		t.Fatalf("expected event round id to be %d, got %d", eventRoundId1, eventRoundUsers[1].EventRoundID)
	}
	if eventRoundUsers[0].Result != 2 {
		t.Fatalf("expected result to be 2, got %d", eventRoundUsers[0].Result)
	}
	if eventRoundUsers[1].Result != 1 {
		t.Fatalf("expected result to be 1, got %d", eventRoundUsers[1].Result)
	}
	if eventRoundUsers[0].Ranking != 1 {
		t.Fatalf("expected ranking to be 1, got %d", eventRoundUsers[0].Ranking)
	}
	if eventRoundUsers[1].Ranking != 2 {
		t.Fatalf("expected ranking to be 2, got %d", eventRoundUsers[1].Ranking)
	}
}

func Test_UpdateEventUser_ByEventUserId_EventUserUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event52",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  31,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	_, err = q.UpdateEventUser(context.Background(), UpdateEventUserParams{
		User: GetEventUserParams{
			ID: sql.NullInt64{Int64: eventUserId, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err != nil {
		t.Fatalf("could not update event user: %v", err)
	}
	row, err := q.GetEventUser(context.Background(), GetEventUserParams{
		ID: sql.NullInt64{Int64: eventUserId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event user: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value2"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value2\"}, got %s", row.Data)
	}
}

func Test_UpdateEventUser_ByEventNameAndUserID_EventUserUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event53",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  32,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	_, err = q.UpdateEventUser(context.Background(), UpdateEventUserParams{
		User: GetEventUserParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event53", Valid: true},
			},
			UserID: sql.NullInt64{Int64: 32, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err != nil {
		t.Fatalf("could not update event user: %v", err)
	}
	row, err := q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event53", Valid: true},
		},
		UserID: sql.NullInt64{Int64: 32, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event user: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value2"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value2\"}, got %s", row.Data)
	}
}

func Test_UpdateEventUser_ByEventIDAndUserID_EventUserUpdated(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event54",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  33,
		Data:    json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	_, err = q.UpdateEventUser(context.Background(), UpdateEventUserParams{
		User: GetEventUserParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: eventId, Valid: true},
			},
			UserID: sql.NullInt64{Int64: 33, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err != nil {
		t.Fatalf("could not update event user: %v", err)
	}
	row, err := q.GetEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
		UserID: sql.NullInt64{Int64: 33, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get event user: %v", err)
	}
	if !reflect.DeepEqual(row.Data, json.RawMessage(`{"key": "value2"}`)) {
		t.Fatalf("expected data to be {\"key\": \"value2\"}, got %s", row.Data)
	}
}

func Test_UpdateEventUser_ByEventIDAndUserIDEventDoesNotExist_EventUserNotUpdated(t *testing.T) {
	q := New(db)
	_, err := q.UpdateEventUser(context.Background(), UpdateEventUserParams{
		User: GetEventUserParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: 99999, Valid: true},
			},
			UserID: sql.NullInt64{Int64: 34, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_UpdateEventUser_ByEventNameAndUserIDEventDoesNotExist_EventUserNotUpdated(t *testing.T) {
	q := New(db)
	_, err := q.UpdateEventUser(context.Background(), UpdateEventUserParams{
		User: GetEventUserParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event55", Valid: true},
			},
			UserID: sql.NullInt64{Int64: 35, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_UpdateEventUser_ByEventUserIdEventDoesNotExist_EventUserNotUpdated(t *testing.T) {
	q := New(db)
	_, err := q.UpdateEventUser(context.Background(), UpdateEventUserParams{
		User: GetEventUserParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventUser_ByEventUserId_EventUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event56",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  36,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.DeleteEventUser(context.Background(), GetEventUserParams{
		ID: sql.NullInt64{Int64: eventUserId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event user: %v", err)
	}
	_, err = q.GetEventUser(context.Background(), GetEventUserParams{
		ID: sql.NullInt64{Int64: eventUserId, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventUser_ByEventNameAndUserID_EventUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event57",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  37,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	result, err = q.DeleteEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event57", Valid: true},
		},
		UserID: sql.NullInt64{Int64: 37, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteEventUser_ByEventIDAndUserID_EventUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event58",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	_, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  38,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	result, err = q.DeleteEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: eventId, Valid: true},
		},
		UserID: sql.NullInt64{Int64: 38, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteEventUser_ByEventIDAndUserIDEventDoesNotExist_EventUserNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		UserID: sql.NullInt64{Int64: 39, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteEventUser_ByEventNameAndUserIDEventDoesNotExist_EventUserNotDeleted(t *testing.T) {
	q := New(db)
	_, err := q.DeleteEventUser(context.Background(), GetEventUserParams{
		Event: GetEventParams{
			Name: sql.NullString{String: "event59", Valid: true},
		},
		UserID: sql.NullInt64{Int64: 40, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_DeleteEventUser_ByEventUserIdEventDoesNotExist_EventUserNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteEventUser(context.Background(), GetEventUserParams{
		ID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteEventRoundUser_ByEventRoundUserId_EventRoundUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event60",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  41,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round45",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	result, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		ID: sql.NullInt64{Int64: eventRoundUserId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventRoundUserIdEventDoesNotExist_EventRoundUserNotDeleted(t *testing.T) {
	q := New(db)
	_, err := q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		ID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventUserIdAndRoundName_EventRoundUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event61",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  42,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round46",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	result, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			ID: sql.NullInt64{Int64: eventUserId, Valid: true},
		},
		Round: sql.NullString{String: "round46", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventUserIdAndRoundNameEventDoesNotExist_EventRoundUserNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			ID: sql.NullInt64{Int64: 99999, Valid: true},
		},
		Round: sql.NullString{String: "round47", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteEventRoundUser_ByEventUserIdCurrentRound_EventRoundUserNotDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event62",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  43,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round48",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	_, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			ID: sql.NullInt64{Int64: eventUserId, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventIdAndUserIdAndRoundName_EventRoundUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event63",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  44,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round49",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	result, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: eventId, Valid: true},
			},
			UserID: sql.NullInt64{Int64: 44, Valid: true},
		},
		Round: sql.NullString{String: "round49", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventIdAndUserIdAndCurrentRound_EventRoundUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event64",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  44,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round50",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	result, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			Event: GetEventParams{
				ID: sql.NullInt64{Int64: eventId, Valid: true},
			},
			UserID: sql.NullInt64{Int64: 44, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventNameAndUserIdAndRoundName_EventRoundUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event65",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  44,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round51",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	result, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event65", Valid: true},
			},
			UserID: sql.NullInt64{Int64: 44, Valid: true},
		},
		Round: sql.NullString{String: "round50", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_DeleteEventRoundUser_ByEventNameAndUserIdAndCurrentRound_EventRoundUserDeleted(t *testing.T) {
	q := New(db)
	startedAt := time.Now().Round(time.Minute)
	result, err := q.CreateEvent(context.Background(), CreateEventParams{
		Name:      "event66",
		Data:      json.RawMessage(`{}`),
		StartedAt: startedAt,
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	eventId, err := result.LastInsertId()
	result, err = q.CreateOrUpdateEventUser(context.Background(), CreateOrUpdateEventUserParams{
		EventID: uint64(eventId),
		UserID:  44,
		Data:    json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create or update event user: %v", err)
	}
	eventUserId, err := result.LastInsertId()
	result, err = q.CreateEventRound(context.Background(), CreateEventRoundParams{
		EventID: uint64(eventId),
		Name:    "round52",
		Data:    json.RawMessage(`{}`),
		Scoring: json.RawMessage(`{"scoring": [4,1]}`),
		EndedAt: startedAt.Add(1 * time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create event round: %v", err)
	}
	eventRoundId, err := result.LastInsertId()
	result, err = q.CreateEventRoundUser(context.Background(), CreateEventRoundUserParams{
		EventRoundID: uint64(eventRoundId),
		EventUserID:  uint64(eventUserId),
		Result:       1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event round user: %v", err)
	}
	eventRoundUserId, err := result.LastInsertId()
	result, err = q.DeleteEventRoundUser(context.Background(), GetEventRoundUserParams{
		EventUser: GetEventUserParams{
			Event: GetEventParams{
				Name: sql.NullString{String: "event66", Valid: true},
			},
			UserID: sql.NullInt64{Int64: 44, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not delete event round user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	_, err = db.QueryContext(context.Background(), `SELECT * FROM event_round_user WHERE id = ?;`, eventRoundUserId)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}
