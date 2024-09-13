package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/mysqlTestServer"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	server, err := mysqlTestServer.GetServer()
	if err != nil {
		log.Fatalf("could not run mysql test server: %v", err)
	}
	defer server.Close()
	db = server.Db
	schema, err := os.ReadFile("../../../migration/tournament.sql")
	if err != nil {
		log.Fatalf("could not read schema file: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("could not execute schema: %v", err)
	}

	m.Run()
}

func Test_CreateTournament_Tournament_TournamentCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTournament(CreateTournamentParams{
		Name:                "test",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              1,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateTournament_TournamentExists_TournamentNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTournament(CreateTournamentParams{
		Name:                "test1",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              2,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(CreateTournamentParams{
		Name:                "test1",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              2,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
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

func Test_GetTournament_

// func Test_CreateRecord_Record_RecordCreated(t *testing.T) {
// 	q := New(db)
// 	result, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 1,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		t.Fatalf("could not get rows affected: %v", err)
// 	}
// 	if rowsAffected != 1 {
// 		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
// 	}
// }

// func Test_CreateRecord_RecordExists_RecordNotCreated(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 2,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 2,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}
// 	mysqlErr, ok := err.(*mysql.MySQLError)
// 	if !ok {
// 		t.Fatalf("expected mysql error, got %v", err)
// 	}
// 	if mysqlErr.Number != errorcode.MySQLErrorCodeDuplicateEntry {
// 		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
// 	}
// }

// func Test_GetRecord_ById_Record(t *testing.T) {
// 	q := New(db)
// 	result, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 3,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		t.Fatalf("could not get last insert id: %v", err)
// 	}
// 	record, err := q.GetRecord(context.Background(), GetRecordParams{
// 		Id: sql.NullInt64{Int64: id, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get record: %v", err)
// 	}
// 	if record.ID != uint64(id) {
// 		t.Fatalf("expected record id 1, got %d", record.ID)
// 	}
// 	if record.Name != "test" {
// 		t.Fatalf("expected record name test, got %s", record.Name)
// 	}
// 	if record.UserID != 3 {
// 		t.Fatalf("expected record user id 1, got %d", record.UserID)
// 	}
// 	if record.Record != 1 {
// 		t.Fatalf("expected record 1, got %d", record.Record)
// 	}
// 	if string(record.Data) != `{"key": "value"}` {
// 		t.Fatalf("expected record data {\"key\": \"value\"}, got %s", record.Data)
// 	}
// }

// func Test_GetRecord_ByNameAndUserId_Record(t *testing.T) {
// 	q := New(db)
// 	result, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 4,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		t.Fatalf("could not get last insert id: %v", err)
// 	}
// 	record, err := q.GetRecord(context.Background(), GetRecordParams{
// 		Name:   sql.NullString{String: "test", Valid: true},
// 		UserID: sql.NullInt64{Int64: 4, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get record: %v", err)
// 	}
// 	if record.ID != uint64(id) {
// 		t.Fatalf("expected record id 1, got %d", record.ID)
// 	}
// 	if record.Name != "test" {
// 		t.Fatalf("expected record name test, got %s", record.Name)
// 	}
// 	if record.UserID != 4 {
// 		t.Fatalf("expected record user id 4, got %d", record.UserID)
// 	}
// 	if record.Record != 1 {
// 		t.Fatalf("expected record 1, got %d", record.Record)
// 	}
// 	if string(record.Data) != `{"key": "value"}` {
// 		t.Fatalf("expected record data {\"key\": \"value\"}, got %s", record.Data)
// 	}
// }

// func Test_GetRecord_RecordDoesNotExist_Error(t *testing.T) {
// 	q := New(db)
// 	_, err := q.GetRecord(context.Background(), GetRecordParams{
// 		Id: sql.NullInt64{Int64: 999999, Valid: true},
// 	})
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}
// 	if err != sql.ErrNoRows {
// 		t.Fatalf("expected sql.ErrNoRows, got %v", err)
// 	}
// }

// func Test_GetRecords_NoNameNoUserId_AllRecords(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 5,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test1",
// 		UserID: 6,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	records, err := q.GetRecords(context.Background(), GetRecordsParams{
// 		Limit:  2,
// 		Offset: 0,
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get records: %v", err)
// 	}
// 	if len(records) != 2 {
// 		t.Fatalf("expected 2 records, got %d", len(records))
// 	}
// }

// func Test_GetRecords_NameNoUserId_Records(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "GetRecords_NameNoUserId_Records",
// 		UserID: 7,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "GetRecords_NameNoUserId_Records",
// 		UserID: 8,
// 		Record: 100,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test1",
// 		UserID: 9,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	records, err := q.GetRecords(context.Background(), GetRecordsParams{
// 		Name:   sql.NullString{String: "GetRecords_NameNoUserId_Records", Valid: true},
// 		Limit:  3,
// 		Offset: 0,
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get records: %v", err)
// 	}
// 	if len(records) != 2 {
// 		t.Fatalf("expected 2 records, got %d", len(records))
// 	}
// 	if records[0].Name != "GetRecords_NameNoUserId_Records" || records[0].UserID != 7 {
// 		t.Fatalf("expected record name GetRecords_NameNoUserId_Records and user id 7, got %s and %d", records[0].Name, records[0].UserID)
// 	}
// 	if records[1].Name != "GetRecords_NameNoUserId_Records" || records[1].UserID != 8 {
// 		t.Fatalf("expected record name GetRecords_NameNoUserId_Records and user id 8, got %s and %d", records[1].Name, records[1].UserID)
// 	}
// 	if records[0].Ranking != 1 {
// 		t.Fatalf("expected first record ranking 1, got %d", records[0].Ranking)
// 	}
// 	if records[1].Ranking != 2 {
// 		t.Fatalf("expected second record ranking 2, got %d", records[1].Ranking)
// 	}
// }

// func Test_GetRecords_NoNameUserId_Records(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 99999,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test1",
// 		UserID: 99999,
// 		Record: 2,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test1",
// 		UserID: 11,
// 		Record: 3,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	records, err := q.GetRecords(context.Background(), GetRecordsParams{
// 		UserId: sql.NullInt64{Int64: 99999, Valid: true},
// 		Limit:  3,
// 		Offset: 0,
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get records: %v", err)
// 	}
// 	if len(records) != 2 {
// 		t.Fatalf("expected 2 records, got %d", len(records))
// 	}
// 	if records[0].Name != "test" || records[0].UserID != 99999 {
// 		t.Fatalf("expected record name test and user id 99999, got %s and %d", records[0].Name, records[0].UserID)
// 	}
// 	if records[1].Name != "test1" || records[1].UserID != 99999 {
// 		t.Fatalf("expected record name test1 and user id 99999, got %s and %d", records[1].Name, records[1].UserID)
// 	}
// }

// func Test_UpdateRecord_UpdateDataById_RecordUpdated(t *testing.T) {
// 	q := New(db)
// 	result, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 12,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		t.Fatalf("could not get last insert id: %v", err)
// 	}
// 	_, err = q.UpdateRecord(context.Background(), UpdateRecordParams{
// 		GetRecordParams: GetRecordParams{Id: sql.NullInt64{Int64: id, Valid: true}},
// 		Data:            json.RawMessage(`{"key": "value1"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not update record: %v", err)
// 	}
// 	record, err := q.GetRecord(context.Background(), GetRecordParams{Id: sql.NullInt64{Int64: id, Valid: true}})
// 	if err != nil {
// 		t.Fatalf("could not get record: %v", err)
// 	}
// 	if string(record.Data) != `{"key": "value1"}` {
// 		t.Fatalf("expected record data {\"key\": \"value1\"}, got %s", record.Data)
// 	}
// }

// func Test_UpdateRecord_UpdateDataByNameUserId_RecordUpdated(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 13,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.UpdateRecord(context.Background(), UpdateRecordParams{
// 		GetRecordParams: GetRecordParams{
// 			Name:   sql.NullString{String: "test", Valid: true},
// 			UserID: sql.NullInt64{Int64: 13, Valid: true},
// 		},
// 		Data: json.RawMessage(`{"key": "value1"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not update record: %v", err)
// 	}
// 	record, err := q.GetRecord(context.Background(), GetRecordParams{
// 		Name:   sql.NullString{String: "test", Valid: true},
// 		UserID: sql.NullInt64{Int64: 13, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get record: %v", err)
// 	}
// 	if string(record.Data) != `{"key": "value1"}` {
// 		t.Fatalf("expected record data {\"key\": \"value1\"}, got %s", record.Data)
// 	}
// }

// func Test_UpdateRecord_UpdateRecordById_RecordUpdated(t *testing.T) {
// 	q := New(db)
// 	result, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 14,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		t.Fatalf("could not get last insert id: %v", err)
// 	}
// 	_, err = q.UpdateRecord(context.Background(), UpdateRecordParams{
// 		GetRecordParams: GetRecordParams{Id: sql.NullInt64{Int64: id, Valid: true}},
// 		Record:          sql.NullInt64{Int64: 2, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not update record: %v", err)
// 	}
// 	record, err := q.GetRecord(context.Background(), GetRecordParams{Id: sql.NullInt64{Int64: id, Valid: true}})
// 	if err != nil {
// 		t.Fatalf("could not get record: %v", err)
// 	}
// 	if record.Record != 2 {
// 		t.Fatalf("expected record 2, got %d", record.Record)
// 	}
// }

// func Test_UpdateRecord_UpdateRecordByNameUserId_RecordUpdated(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 15,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.UpdateRecord(context.Background(), UpdateRecordParams{
// 		GetRecordParams: GetRecordParams{
// 			Name:   sql.NullString{String: "test", Valid: true},
// 			UserID: sql.NullInt64{Int64: 15, Valid: true},
// 		},
// 		Record: sql.NullInt64{Int64: 2, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not update record: %v", err)
// 	}
// 	record, err := q.GetRecord(context.Background(), GetRecordParams{
// 		Name:   sql.NullString{String: "test", Valid: true},
// 		UserID: sql.NullInt64{Int64: 15, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get record: %v", err)
// 	}
// 	if record.Record != 2 {
// 		t.Fatalf("expected record 2, got %d", record.Record)
// 	}
// }

// func Test_UpdateRecord_RecordDoesNotExist_Error(t *testing.T) {
// 	q := New(db)
// 	result, err := q.UpdateRecord(context.Background(), UpdateRecordParams{
// 		GetRecordParams: GetRecordParams{Id: sql.NullInt64{Int64: 999999, Valid: true}},
// 		Data:            json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not update record: %v", err)
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		t.Fatalf("could not get rows affected: %v", err)
// 	}
// 	if rowsAffected != 0 {
// 		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
// 	}

// }

// func Test_DeleteRecord_DeleteById_RecordDeleted(t *testing.T) {
// 	q := New(db)
// 	result, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 16,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		t.Fatalf("could not get last insert id: %v", err)
// 	}
// 	_, err = q.DeleteRecord(context.Background(), GetRecordParams{Id: sql.NullInt64{Int64: id, Valid: true}})
// 	if err != nil {
// 		t.Fatalf("could not delete record: %v", err)
// 	}
// 	_, err = q.GetRecord(context.Background(), GetRecordParams{Id: sql.NullInt64{Int64: id, Valid: true}})
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}
// 	if err != sql.ErrNoRows {
// 		t.Fatalf("expected sql.ErrNoRows, got %v", err)
// 	}
// }

// func Test_DeleteRecord_DeleteByNameUserId_RecordDeleted(t *testing.T) {
// 	q := New(db)
// 	_, err := q.CreateRecord(context.Background(), CreateRecordParams{
// 		Name:   "test",
// 		UserID: 17,
// 		Record: 1,
// 		Data:   json.RawMessage(`{"key": "value"}`),
// 	})
// 	if err != nil {
// 		t.Fatalf("could not create record: %v", err)
// 	}
// 	_, err = q.DeleteRecord(context.Background(), GetRecordParams{
// 		Name:   sql.NullString{String: "test", Valid: true},
// 		UserID: sql.NullInt64{Int64: 17, Valid: true},
// 	})
// 	if err != nil {
// 		t.Fatalf("could not delete record: %v", err)
// 	}
// 	_, err = q.GetRecord(context.Background(), GetRecordParams{
// 		Name:   sql.NullString{String: "test", Valid: true},
// 		UserID: sql.NullInt64{Int64: 17, Valid: true},
// 	})
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}
// 	if err != sql.ErrNoRows {
// 		t.Fatalf("expected sql.ErrNoRows, got %v", err)
// 	}
// }

// func Test_DeleteRecord_RecordDoesNotExist_Error(t *testing.T) {
// 	q := New(db)
// 	result, err := q.DeleteRecord(context.Background(), GetRecordParams{Id: sql.NullInt64{Int64: 999999, Valid: true}})
// 	if err != nil {
// 		t.Fatalf("could not delete record: %v", err)
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		t.Fatalf("could not get rows affected: %v", err)
// 	}
// 	if rowsAffected != 0 {
// 		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
// 	}
// }
