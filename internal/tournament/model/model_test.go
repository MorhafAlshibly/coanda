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
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
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
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test1",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              2,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test1",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              2,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
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

func Test_CreateTournament_SameNameDifferentUser_TournamentCreated(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test2",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              3,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test2",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              4,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
}

func Test_CreateTournament_SameNameSameUserDifferentInterval_TournamentCreated(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test3",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              5,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test3",
		TournamentInterval:  TournamentTournamentIntervalWeekly,
		UserID:              5,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
}

func Test_CreateTournament_SameNameSameUserSameIntervalDifferentStartedAt_TournamentCreated(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test4",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              6,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test4",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              6,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
}

func Test_GetTournament_ById_Tournament(t *testing.T) {
	q := New(db)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test5",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              3,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if tournament.ID != uint64(id) {
		t.Fatalf("expected tournament id %d, got %d", id, tournament.ID)
	}
	if tournament.Name != "test5" {
		t.Fatalf("expected tournament name test5, got %s", tournament.Name)
	}
	if tournament.UserID != 3 {
		t.Fatalf("expected tournament user id 1, got %d", tournament.UserID)
	}
	if tournament.Score != 1 {
		t.Fatalf("expected tournament score 1, got %d", tournament.Score)
	}
	if string(tournament.Data) != `{"key": "value"}` {
		t.Fatalf("expected tournament data {\"key\": \"value\"}, got %s", tournament.Data)
	}
}

func Test_GetTournament_ByNameIntervalUserIdStartedAt_Tournament(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test6",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              4,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		NameIntervalUserIDStartedAt: NullNameIntervalUserIDStartedAt{
			Name:                "test6",
			TournamentInterval:  TournamentTournamentIntervalDaily,
			UserID:              4,
			TournamentStartedAt: tournamentStartedAt,
			Valid:               true,
		},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if tournament.ID != uint64(id) {
		t.Fatalf("expected tournament id %d, got %d", id, tournament.ID)
	}
	if tournament.Name != "test6" {
		t.Fatalf("expected tournament name test6, got %s", tournament.Name)
	}
	if tournament.UserID != 4 {
		t.Fatalf("expected tournament user id 4, got %d", tournament.UserID)
	}
	if tournament.Score != 1 {
		t.Fatalf("expected tournament score 1, got %d", tournament.Score)
	}
	if string(tournament.Data) != `{"key": "value"}` {
		t.Fatalf("expected tournament data {\"key\": \"value\"}, got %s", tournament.Data)
	}
	if !tournament.TournamentStartedAt.Equal(tournamentStartedAt) {
		t.Fatalf("expected tournament started at %v, got %v", tournamentStartedAt, tournament.TournamentStartedAt)
	}
}

func Test_GetTournament_TournamentDoesNotExist_Error(t *testing.T) {
	q := New(db)
	_, err := q.GetTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: 999999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func Test_GetTournament_ById_TournamentUserRankedSecond(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test7",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              5,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test7",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              6,
		Score:               2,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if tournament.Ranking != 2 {
		t.Fatalf("expected tournament ranking 2, got %d", tournament.Ranking)
	}
}

func Test_GetTournaments_NoNameNoUserId_AllTournaments(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test8",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              7,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test8",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              8,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.GetTournaments(context.Background(), GetTournamentsParams{
		TournamentInterval:  TournamentTournamentIntervalDaily,
		TournamentStartedAt: tournamentStartedAt,
		Limit:               2,
		Offset:              0,
	})
	if err != nil {
		t.Fatalf("could not get tournaments: %v", err)
	}
}

func Test_GetTournaments_NameNoUserId_Tournaments(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test9",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              9,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test9",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              10,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test10",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              11,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	tournaments, err := q.GetTournaments(context.Background(), GetTournamentsParams{
		Name:                sql.NullString{String: "test9", Valid: true},
		TournamentInterval:  TournamentTournamentIntervalDaily,
		TournamentStartedAt: tournamentStartedAt,
		Limit:               3,
		Offset:              0,
	})
	if err != nil {
		t.Fatalf("could not get tournaments: %v", err)
	}
	if len(tournaments) != 2 {
		t.Fatalf("expected 2 tournaments, got %d", len(tournaments))
	}
}

func Test_GetTournaments_NoNameUserId_Tournaments(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test11",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              12,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test12",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              12,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	_, err = q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test12",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              11,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	tournaments, err := q.GetTournaments(context.Background(), GetTournamentsParams{
		UserID:              sql.NullInt64{Int64: 12, Valid: true},
		TournamentInterval:  TournamentTournamentIntervalDaily,
		TournamentStartedAt: tournamentStartedAt,
		Limit:               3,
		Offset:              0,
	})
	if err != nil {
		t.Fatalf("could not get tournaments: %v", err)
	}
	if len(tournaments) != 2 {
		t.Fatalf("expected 2 tournaments, got %d", len(tournaments))
	}
}

func Test_DeleteTournament_ById_TournamentDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test13",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              13,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.DeleteTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteTournament_ById_TournamentDoesNotExist_Error(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: 999999, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteTournament_ByNameIntervalUserIdStartedAt_TournamentDeleted(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test14",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              14,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	result, err = q.DeleteTournament(context.Background(), GetTournamentParams{
		NameIntervalUserIDStartedAt: NullNameIntervalUserIDStartedAt{
			Name:                "test14",
			TournamentInterval:  TournamentTournamentIntervalDaily,
			UserID:              14,
			TournamentStartedAt: tournamentStartedAt,
			Valid:               true,
		},
	})
	if err != nil {
		t.Fatalf("could not delete tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteTournament_ByNameIntervalUserIdStartedAt_TournamentDoesNotExist_Error(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTournament(context.Background(), GetTournamentParams{
		NameIntervalUserIDStartedAt: NullNameIntervalUserIDStartedAt{
			Name:                "test15",
			TournamentInterval:  TournamentTournamentIntervalDaily,
			UserID:              15,
			TournamentStartedAt: time.Now(),
			Valid:               true,
		},
	})
	if err != nil {
		t.Fatalf("could not delete tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateTournament_UpdateDataById_TournamentUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test16",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              16,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.UpdateTournament(context.Background(), UpdateTournamentParams{
		Tournament: GetTournamentParams{
			ID: sql.NullInt64{Int64: id, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err != nil {
		t.Fatalf("could not update tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if string(tournament.Data) != `{"key": "value2"}` {
		t.Fatalf("expected tournament data {\"key\": \"value2\"}, got %s", tournament.Data)
	}
}

func Test_UpdateTournament_UpdateScoreById_TournamentUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test17",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              17,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.UpdateTournament(context.Background(), UpdateTournamentParams{
		Tournament: GetTournamentParams{
			ID: sql.NullInt64{Int64: id, Valid: true},
		},
		Score:          sql.NullInt64{Int64: 2, Valid: true},
		IncrementScore: false,
	})
	if err != nil {
		t.Fatalf("could not update tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if tournament.Score != 2 {
		t.Fatalf("expected tournament score 2, got %d", tournament.Score)
	}
}

func Test_UpdateTournament_IncrementScoreById_TournamentUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test18",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              18,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.UpdateTournament(context.Background(), UpdateTournamentParams{
		Tournament: GetTournamentParams{
			ID: sql.NullInt64{Int64: id, Valid: true},
		},
		Score:          sql.NullInt64{Int64: 2, Valid: true},
		IncrementScore: true,
	})
	if err != nil {
		t.Fatalf("could not update tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		ID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if tournament.Score != 3 {
		t.Fatalf("expected tournament score 3, got %d", tournament.Score)
	}
}

func Test_UpdateTournament_ById_TournamentDoesNotExist_Error(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTournament(context.Background(), UpdateTournamentParams{
		Tournament: GetTournamentParams{
			ID: sql.NullInt64{Int64: 999999, Valid: true},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateTournament_ByNameIntervalUserIdStartedAt_TournamentUpdated(t *testing.T) {
	q := New(db)
	tournamentStartedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := q.CreateTournament(context.Background(), CreateTournamentParams{
		Name:                "test19",
		TournamentInterval:  TournamentTournamentIntervalDaily,
		UserID:              19,
		Score:               1,
		Data:                json.RawMessage(`{"key": "value"}`),
		TournamentStartedAt: tournamentStartedAt,
	})
	if err != nil {
		t.Fatalf("could not create tournament: %v", err)
	}
	result, err = q.UpdateTournament(context.Background(), UpdateTournamentParams{
		Tournament: GetTournamentParams{
			NameIntervalUserIDStartedAt: NullNameIntervalUserIDStartedAt{
				Name:                "test19",
				TournamentInterval:  TournamentTournamentIntervalDaily,
				UserID:              19,
				TournamentStartedAt: tournamentStartedAt,
				Valid:               true,
			},
		},
		Data: json.RawMessage(`{"key": "value2"}`),
	})
	if err != nil {
		t.Fatalf("could not update tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	tournament, err := q.GetTournament(context.Background(), GetTournamentParams{
		NameIntervalUserIDStartedAt: NullNameIntervalUserIDStartedAt{
			Name:                "test19",
			TournamentInterval:  TournamentTournamentIntervalDaily,
			UserID:              19,
			TournamentStartedAt: tournamentStartedAt,
			Valid:               true,
		},
	})
	if err != nil {
		t.Fatalf("could not get tournament: %v", err)
	}
	if string(tournament.Data) != `{"key": "value2"}` {
		t.Fatalf("expected tournament data {\"key\": \"value2\"}, got %s", tournament.Data)
	}
}

func Test_UpdateTournament_ByNameIntervalUserIdStartedAt_TournamentDoesNotExist_Error(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTournament(context.Background(), UpdateTournamentParams{
		Tournament: GetTournamentParams{
			NameIntervalUserIDStartedAt: NullNameIntervalUserIDStartedAt{
				Name:                "test20",
				TournamentInterval:  TournamentTournamentIntervalDaily,
				UserID:              20,
				TournamentStartedAt: time.Now(),
				Valid:               true,
			},
		},
		Data: json.RawMessage(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("could not update tournament: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}
