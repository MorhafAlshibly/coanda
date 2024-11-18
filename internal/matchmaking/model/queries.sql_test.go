package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/MorhafAlshibly/coanda/pkg/mysqlTestServer"
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
