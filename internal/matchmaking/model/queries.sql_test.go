package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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

func Test_CreateMatchmakingTicketUser_MatchmakingTicketUser_MatchmakingTicketUserCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 4,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	userId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingTicket(context.Background(), CreateMatchmakingTicketParams{
		Data:      json.RawMessage(`{}`),
		EloWindow: 0,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	ticketId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingTicketUser(context.Background(), CreateMatchmakingTicketUserParams{
		MatchmakingUserID:   uint64(userId),
		MatchmakingTicketID: uint64(ticketId),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateMatchmakingTicketArena_MatchmakingTicketArena_MatchmakingTicketArenaCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena3",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	arenaId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingTicket(context.Background(), CreateMatchmakingTicketParams{
		Data:      json.RawMessage(`{}`),
		EloWindow: 0,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	ticketId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingTicketArena(context.Background(), CreateMatchmakingTicketArenaParams{
		MatchmakingTicketID: uint64(ticketId),
		MatchmakingArenaID:  uint64(arenaId),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket arena: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_UpdateMatchmakingUserByClientUserId_ValidClientUserId_UserUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 5,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	result, err := q.UpdateMatchmakingUserByClientUserId(context.Background(), UpdateMatchmakingUserByClientUserIdParams{
		ClientUserID: 5,
		Elo:          2000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not update matchmaking user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	var checkElo int
	err = db.QueryRow("SELECT elo FROM matchmaking_user WHERE client_user_id = ?", 5).Scan(&checkElo)
	if err != nil {
		t.Fatalf("could not scan row: %v", err)
	}
	if checkElo != 2000 {
		t.Fatalf("expected elo 2000, got %d", checkElo)
	}
}

func Test_UpdateMatchmakingUserByClientUserId_InvalidClientUserId_UserNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateMatchmakingUserByClientUserId(context.Background(), UpdateMatchmakingUserByClientUserIdParams{
		ClientUserID: 999999999,
		Elo:          2000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not update matchmaking user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetArena_ByID_ArenaExists(t *testing.T) {
	q := New(db)
	result, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena4",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	arenaId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	arena, err := q.GetArena(context.Background(), ArenaParams{
		ID: conversion.Int64ToSqlNullInt64(&arenaId),
	})
	if err != nil {
		t.Fatalf("could not get arena: %v", err)
	}
	if arena.ID != uint64(arenaId) {
		t.Fatalf("expected arena id %d, got %d", arenaId, arena.ID)
	}
	if arena.Name != "arena4" {
		t.Fatalf("expected arena name arena4, got %s", arena.Name)
	}
	if arena.MinPlayers != 5 {
		t.Fatalf("expected min players 5, got %d", arena.MinPlayers)
	}
	if arena.MaxPlayersPerTicket != 7 {
		t.Fatalf("expected max players per ticket 7, got %d", arena.MaxPlayersPerTicket)
	}
	if arena.MaxPlayers != 20 {
		t.Fatalf("expected max players 20, got %d", arena.MaxPlayers)
	}
	if arena.Data == nil {
		t.Fatalf("expected non-nil data")
	}
	if arena.CreatedAt.IsZero() {
		t.Fatalf("expected non-zero created at")
	}
	if arena.UpdatedAt.IsZero() {
		t.Fatalf("expected non-zero updated at")
	}
}

func Test_GetArena_ByID_ArenaDoesntExist(t *testing.T) {
	q := New(db)
	_, err := q.GetArena(context.Background(), ArenaParams{
		ID: sql.NullInt64{
			Int64: 999999999,
			Valid: true,
		},
	})
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetArena_ByName_ArenaExists(t *testing.T) {
	q := New(db)
	result, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena5",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	arenaId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	arena, err := q.GetArena(context.Background(), ArenaParams{
		Name: sql.NullString{
			String: "arena5",
			Valid:  true,
		},
	})
	if err != nil {
		t.Fatalf("could not get arena: %v", err)
	}
	if arena.ID != uint64(arenaId) {
		t.Fatalf("expected arena id %d, got %d", arenaId, arena.ID)
	}
	if arena.Name != "arena5" {
		t.Fatalf("expected arena name arena5, got %s", arena.Name)
	}
	if arena.MinPlayers != 5 {
		t.Fatalf("expected min players 5, got %d", arena.MinPlayers)
	}
	if arena.MaxPlayersPerTicket != 7 {
		t.Fatalf("expected max players per ticket 7, got %d", arena.MaxPlayersPerTicket)
	}
	if arena.MaxPlayers != 20 {
		t.Fatalf("expected max players 20, got %d", arena.MaxPlayers)
	}
	if arena.Data == nil {
		t.Fatalf("expected non-nil data")
	}
	if arena.CreatedAt.IsZero() {
		t.Fatalf("expected non-zero created at")
	}
	if arena.UpdatedAt.IsZero() {
		t.Fatalf("expected non-zero updated at")
	}
}

func Test_GetArena_ByName_ArenaDoesntExist(t *testing.T) {
	q := New(db)
	_, err := q.GetArena(context.Background(), ArenaParams{
		Name: sql.NullString{
			String: "arena999999",
			Valid:  true,
		},
	})
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

// TODO: UpdateArena tests go here

// TODO: GetMatchmakingUser tests go here

// TODO: UpdateMatchmakingUser tests go here

func Test_GetMatchmakingTicket_ByID_MatchmakingTicketExists(t *testing.T) {
	q := New(db)
	result, err := q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena6",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	arenaId1, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateArena(context.Background(), CreateArenaParams{
		Name:                "arena7",
		MinPlayers:          5,
		MaxPlayersPerTicket: 7,
		MaxPlayers:          20,
		Data:                json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create arena: %v", err)
	}
	arenaId2, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 6,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	userId1, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
		ClientUserID: 7,
		Elo:          1000,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking user: %v", err)
	}
	userId2, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingTicket(context.Background(), CreateMatchmakingTicketParams{
		Data:      json.RawMessage(`{}`),
		EloWindow: 0,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	ticketId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingTicketUser(context.Background(), CreateMatchmakingTicketUserParams{
		MatchmakingUserID:   uint64(userId1),
		MatchmakingTicketID: uint64(ticketId),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket user: %v", err)
	}
	result, err = q.CreateMatchmakingTicketUser(context.Background(), CreateMatchmakingTicketUserParams{
		MatchmakingUserID:   uint64(userId2),
		MatchmakingTicketID: uint64(ticketId),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket user: %v", err)
	}
	result, err = q.CreateMatchmakingTicketArena(context.Background(), CreateMatchmakingTicketArenaParams{
		MatchmakingTicketID: uint64(ticketId),
		MatchmakingArenaID:  uint64(arenaId1),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket arena: %v", err)
	}
	result, err = q.CreateMatchmakingTicketArena(context.Background(), CreateMatchmakingTicketArenaParams{
		MatchmakingTicketID: uint64(ticketId),
		MatchmakingArenaID:  uint64(arenaId2),
	})
	if err != nil {
		t.Fatalf("could not create matchmaking ticket arena: %v", err)
	}
	matchmakingTicketRows, err := q.GetMatchmakingTicket(context.Background(), GetMatchmakingTicketParams{
		MatchmakingTicket: MatchmakingTicketParams{
			ID:                        conversion.Int64ToSqlNullInt64(&ticketId),
			GetByIDRegardlessOfStatus: true,
		},
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking ticket: %v", err)
	}
	if matchmakingTicketRows[0].TicketID != uint64(ticketId) {
		t.Fatalf("expected ticket id %d, got %d", ticketId, matchmakingTicketRows[0].TicketID)
	}
	if matchmakingTicketRows[0].MatchmakingUserID != uint64(userId1) {
		t.Fatalf("expected user id %d, got %d", userId1, matchmakingTicketRows[0].MatchmakingUserID)
	}
	if matchmakingTicketRows[0].ArenaID != uint64(arenaId1) {
		t.Fatalf("expected arena id %d, got %d", arenaId1, matchmakingTicketRows[0].ArenaID)
	}
	if matchmakingTicketRows[1].MatchmakingUserID != uint64(userId1) {
		t.Fatalf("expected user id %d, got %d", userId1, matchmakingTicketRows[1].MatchmakingUserID)
	}
	if matchmakingTicketRows[1].ArenaID != uint64(arenaId2) {
		t.Fatalf("expected arena id %d, got %d", arenaId2, matchmakingTicketRows[1].ArenaID)
	}
	if matchmakingTicketRows[2].MatchmakingUserID != uint64(userId2) {
		t.Fatalf("expected user id %d, got %d", userId2, matchmakingTicketRows[2].MatchmakingUserID)
	}
	if matchmakingTicketRows[2].ArenaID != uint64(arenaId1) {
		t.Fatalf("expected arena id %d, got %d", arenaId1, matchmakingTicketRows[2].ArenaID)
	}
	if matchmakingTicketRows[3].MatchmakingUserID != uint64(userId2) {
		t.Fatalf("expected user id %d, got %d", userId2, matchmakingTicketRows[3].MatchmakingUserID)
	}
	if matchmakingTicketRows[3].ArenaID != uint64(arenaId2) {
		t.Fatalf("expected arena id %d, got %d", arenaId2, matchmakingTicketRows[3].ArenaID)
	}
}
