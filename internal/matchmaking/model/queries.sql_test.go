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
	schema, err := os.ReadFile("../../../migration/matchmaking.sql")
	if err != nil {
		log.Fatalf("could not read schema file: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("could not execute schema: %v", err)
	}

	m.Run()
}

func Test_CreateArena_Arena_ArenaCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateArena_ArenaExists_ArenaNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena1",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create event: %v", err)
	}
	_, err = q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena1",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if !errorcode.IsDuplicateEntry(mysqlErr, "arena", "name") {
		t.Fatalf("expected duplicate entry error, got %v", err)
	}
}

func Test_GetArenas_OneArena_Arenas(t *testing.T) {
	q := New(db)
	_, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena2",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	arenas, err := q.GetArenas(context.Background(), GetArenasParams{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("could not get arenas: %v", err)
	}
	if len(arenas) != 1 {
		t.Fatalf("expected 1 arena, got %d", len(arenas))
	}
}

func Test_GetArenas_NoArenas_NoArenas(t *testing.T) {
	q := New(db)
	arenas, err := q.GetArenas(context.Background(), GetArenasParams{
		Limit: 0,
	})
	if err != nil {
		t.Fatalf("could not get arenas: %v", err)
	}
	if len(arenas) != 0 {
		t.Fatalf("expected 0 arenas, got %d", len(arenas))
	}
}

func Test_CreateMatchmakingUser_MatchmakingUser_MatchmakingUserCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 1,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateMatchmakingUser_MatchmakingUserExists_MatchmakingUserNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 2,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	_, err = q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 2,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if !errorcode.IsDuplicateEntry(mysqlErr, "matchmaking_user", "client_user_id") {
		t.Fatalf("expected duplicate entry error, got %v", err)
	}
}

func Test_GetMatchmakingUsers_OneMatchmakingUser_MatchmakingUsers(t *testing.T) {
	q := New(db)
	_, err := q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 3,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	matchmakingUsers, err := q.GetMatchmakingUsers(context.Background(), GetMatchmakingUsersParams{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking users: %v", err)
	}
	if len(matchmakingUsers) != 1 {
		t.Fatalf("expected 1 matchmaking user, got %d", len(matchmakingUsers))
	}
}

func Test_GetMatchmakingUsers_NoMatchmakingUsers_NoMatchmakingUsers(t *testing.T) {
	q := New(db)
	matchmakingUsers, err := q.GetMatchmakingUsers(context.Background(), GetMatchmakingUsersParams{
		Limit: 0,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking users: %v", err)
	}
	if len(matchmakingUsers) != 0 {
		t.Fatalf("expected 0 matchmaking users, got %d", len(matchmakingUsers))
	}
}

func Test_CreateMatchmakingTicket_MatchmakingTicket_MatchmakingTicketCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateMatchmakingTicket(context.Background(), CreateMatchmakingTicketParams{
		Data:      json.RawMessage(`{}`),
		EloWindow: 0,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}
