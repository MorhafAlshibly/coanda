package model

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/mysqlTestServer"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var server *mysqlTestServer.Server

func TestMain(m *testing.M) {
	server = mysqlTestServer.NewServer("../../../migration/matchmaking.sql")
	defer server.Close()
	m.Run()
}

func Test_CreateArena_Arena_ArenaCreated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
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
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.CreateMatchmakingTicket(context.Background(), json.RawMessage(`{}`))
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

func Test_AddTicketIdToUser_MatchmakingTicketUser_MatchmakingTicketUserCreated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	result, err := q.CreateMatchmakingTicket(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	ticketId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
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
	result, err = q.AddTicketIDToUser(context.Background(), AddTicketIDToUserParams{
		ID:                  uint64(userId),
		MatchmakingTicketID: conversion.Int64ToSqlNullInt64(&ticketId),
	})
	if err != nil {
		t.Fatalf("could not add ticket id to user: %v", err)
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
	tx := server.Connect(t)
	q := New(tx)
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
	result, err = q.CreateMatchmakingTicket(context.Background(), json.RawMessage(`{}`))
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

func Test_GetArena_ByID_ArenaExists(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
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
	}, nil)
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
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.GetArena(context.Background(), ArenaParams{
		ID: sql.NullInt64{
			Int64: 999999999,
			Valid: true,
		},
	}, nil)
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetArena_ByName_ArenaExists(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
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
	}, nil)
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
	tx := server.Connect(t)
	q := New(tx)
	_, err := q.GetArena(context.Background(), ArenaParams{
		Name: sql.NullString{
			String: "arena999999",
			Valid:  true,
		},
	}, nil)
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

// TODO: UpdateArena tests go here

// TODO: GetMatchmakingUser tests go here

// TODO: UpdateMatchmakingUser tests go here

type TestTicketData struct {
	userId  []int64
	arenaId []int64
}

func createTestTickets(q *Queries) ([]int64, []TestTicketData, error) {
	// Create 4 arenas
	arenaIds := make([]int64, 4)
	for i := 0; i < 4; i++ {
		suffix, err := rand.Int(rand.Reader, big.NewInt(9999999))
		if err != nil {
			return nil, nil, err
		}
		result, err := q.CreateArena(context.Background(), CreateArenaParams{
			Name:                "arenaForTestTickets" + suffix.String(),
			MinPlayers:          5,
			MaxPlayersPerTicket: 7,
			MaxPlayers:          20,
			Data:                json.RawMessage(`{}`),
		})
		if err != nil {
			return nil, nil, err
		}
		arenaId, err := result.LastInsertId()
		if err != nil {
			return nil, nil, err
		}
		arenaIds[i] = arenaId
	}
	// Create 7 matchmaking users
	userIds := make([]int64, 7)
	for i := 0; i < 7; i++ {
		suffix, err := rand.Int(rand.Reader, big.NewInt(9999999))
		if err != nil {
			return nil, nil, err
		}
		result, err := q.CreateMatchmakingUser(context.Background(), CreateMatchmakingUserParams{
			ClientUserID: uint64(suffix.Int64()),
			Elo:          1000,
			Data:         json.RawMessage(`{}`),
		})
		if err != nil {
			return nil, nil, err
		}
		userId, err := result.LastInsertId()
		if err != nil {
			return nil, nil, err
		}
		userIds[i] = userId
	}
	// Create 3 matchmaking tickets
	ticketData := []TestTicketData{
		{
			userId:  []int64{userIds[0], userIds[1]},
			arenaId: []int64{arenaIds[0], arenaIds[1]},
		},
		{
			userId:  []int64{userIds[2], userIds[3]},
			arenaId: []int64{arenaIds[1], arenaIds[2], arenaIds[3]},
		},
		{
			userId:  []int64{userIds[4], userIds[5], userIds[6]},
			arenaId: []int64{arenaIds[3]},
		},
	}
	ticketIds := make([]int64, 3)
	for i, data := range ticketData {
		result, err := q.CreateMatchmakingTicket(context.Background(), json.RawMessage(`{}`))
		if err != nil {
			return nil, nil, err
		}
		ticketId, err := result.LastInsertId()
		if err != nil {
			return nil, nil, err
		}
		ticketIds[i] = ticketId
		for _, userId := range data.userId {
			_, err := q.AddTicketIDToUser(context.Background(), AddTicketIDToUserParams{
				ID:                  uint64(userId),
				MatchmakingTicketID: conversion.Int64ToSqlNullInt64(&ticketId),
			})
			if err != nil {
				return nil, nil, err
			}
		}
		for _, arenaId := range data.arenaId {
			_, err := q.CreateMatchmakingTicketArena(context.Background(), CreateMatchmakingTicketArenaParams{
				MatchmakingTicketID: uint64(ticketId),
				MatchmakingArenaID:  uint64(arenaId),
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return ticketIds, ticketData, nil
}

func Test_GetMatchmakingTicket_ByID_MatchmakingTicketExists(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	matchmakingTicketRows, err := q.GetMatchmakingTicket(context.Background(), GetMatchmakingTicketParams{
		MatchmakingTicket: MatchmakingTicketParams{
			ID: conversion.Int64ToSqlNullInt64(&ticketIds[0]),
		},
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking ticket: %v", err)
	}
	if len(matchmakingTicketRows) != 4 {
		t.Fatalf("expected 4 matchmaking tickets, got %d", len(matchmakingTicketRows))
	}
	if matchmakingTicketRows[0].TicketID != uint64(ticketIds[0]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matchmakingTicketRows[0].TicketID)
	}
	if matchmakingTicketRows[0].MatchmakingUserID != uint64(ticketData[0].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matchmakingTicketRows[0].MatchmakingUserID)
	}
	if matchmakingTicketRows[0].ArenaID != uint64(ticketData[0].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[0], matchmakingTicketRows[0].ArenaID)
	}
	if matchmakingTicketRows[1].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matchmakingTicketRows[1].ArenaID)
	}
	if matchmakingTicketRows[2].MatchmakingUserID != uint64(ticketData[0].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matchmakingTicketRows[2].MatchmakingUserID)
	}
}

func Test_GetMatchmakingTicket_ByMatchmakingUserID_MatchmakingTicketExists(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	matchmakingTicketRows, err := q.GetMatchmakingTicket(context.Background(), GetMatchmakingTicketParams{
		MatchmakingTicket: MatchmakingTicketParams{
			MatchmakingUser: MatchmakingUserParams{
				ID: conversion.Int64ToSqlNullInt64(&ticketData[0].userId[0]),
			},
		},
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking ticket: %v", err)
	}
	if len(matchmakingTicketRows) != 4 {
		t.Fatalf("expected 4 matchmaking tickets, got %d", len(matchmakingTicketRows))
	}
	if matchmakingTicketRows[0].TicketID != uint64(ticketIds[0]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matchmakingTicketRows[0].TicketID)
	}
	if matchmakingTicketRows[0].MatchmakingUserID != uint64(ticketData[0].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matchmakingTicketRows[0].MatchmakingUserID)
	}
	if matchmakingTicketRows[0].ArenaID != uint64(ticketData[0].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[0], matchmakingTicketRows[0].ArenaID)
	}
	if matchmakingTicketRows[1].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matchmakingTicketRows[1].ArenaID)
	}
	if matchmakingTicketRows[2].MatchmakingUserID != uint64(ticketData[0].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matchmakingTicketRows[2].MatchmakingUserID)
	}
}

func Test_GetMatchmakingTicket_ByMatchmakingUserIDWithStatus_MatchmakingTicketExists(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create matchmaking ticket: %v", err)
	}
	matchmakingTicketRows, err := q.GetMatchmakingTicket(context.Background(), GetMatchmakingTicketParams{
		MatchmakingTicket: MatchmakingTicketParams{
			MatchmakingUser: MatchmakingUserParams{
				ID: conversion.Int64ToSqlNullInt64(&ticketData[0].userId[0]),
			},
		},
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking ticket: %v", err)
	}
	if len(matchmakingTicketRows) != 4 {
		t.Fatalf("expected 4 matchmaking tickets, got %d", len(matchmakingTicketRows))
	}
	if matchmakingTicketRows[0].TicketID != uint64(ticketIds[0]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matchmakingTicketRows[0].TicketID)
	}
	if matchmakingTicketRows[0].MatchmakingUserID != uint64(ticketData[0].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matchmakingTicketRows[0].MatchmakingUserID)
	}
	if matchmakingTicketRows[0].ArenaID != uint64(ticketData[0].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[0], matchmakingTicketRows[0].ArenaID)
	}
	if matchmakingTicketRows[1].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matchmakingTicketRows[1].ArenaID)
	}
	if matchmakingTicketRows[2].MatchmakingUserID != uint64(ticketData[0].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matchmakingTicketRows[2].MatchmakingUserID)
	}
}

func Test_GetMatchmakingTicket_ByID_MatchmakingTicketDoesntExist(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	tickets, err := q.GetMatchmakingTicket(context.Background(), GetMatchmakingTicketParams{
		MatchmakingTicket: MatchmakingTicketParams{
			ID: sql.NullInt64{Int64: 999999999, Valid: true},
		},
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking ticket: %v", err)
	}
	if len(tickets) != 0 {
		t.Fatalf("expected 0 tickets, got %d", len(tickets))
	}
}

func Test_GetMatchmakingTickets_NoFilters_TicketsReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	tickets, err := q.GetMatchmakingTickets(context.Background(), GetMatchmakingTicketsParams{
		Limit:      10,
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking tickets: %v", err)
	}
	if len(tickets) != 13 {
		t.Fatalf("expected 13 tickets, got %d", len(tickets))
	}
	if tickets[0].TicketID != uint64(ticketIds[0]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], tickets[0].TicketID)
	}
	if tickets[0].MatchmakingUserID != uint64(ticketData[0].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], tickets[0].MatchmakingUserID)
	}
	if tickets[0].ArenaID != uint64(ticketData[0].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[0], tickets[0].ArenaID)
	}
	if tickets[1].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], tickets[1].ArenaID)
	}
	if tickets[2].MatchmakingUserID != uint64(ticketData[0].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], tickets[2].MatchmakingUserID)
	}
	if tickets[4].TicketID != uint64(ticketIds[1]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[1], tickets[4].TicketID)
	}
	if tickets[4].MatchmakingUserID != uint64(ticketData[1].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[0], tickets[4].MatchmakingUserID)
	}
	if tickets[4].ArenaID != uint64(ticketData[1].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[1].arenaId[0], tickets[4].ArenaID)
	}
	if tickets[5].ArenaID != uint64(ticketData[1].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[1].arenaId[1], tickets[5].ArenaID)
	}
	if tickets[6].ArenaID != uint64(ticketData[1].arenaId[2]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[1].arenaId[2], tickets[6].ArenaID)
	}
	if tickets[7].MatchmakingUserID != uint64(ticketData[1].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[1], tickets[7].MatchmakingUserID)
	}
	if tickets[10].TicketID != uint64(ticketIds[2]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[2], tickets[10].TicketID)
	}
	if tickets[10].MatchmakingUserID != uint64(ticketData[2].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[2].userId[0], tickets[10].MatchmakingUserID)
	}
	if tickets[10].ArenaID != uint64(ticketData[2].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[2].arenaId[0], tickets[10].ArenaID)
	}
	if tickets[11].MatchmakingUserID != uint64(ticketData[2].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[2].userId[1], tickets[11].MatchmakingUserID)
	}
	if tickets[12].MatchmakingUserID != uint64(ticketData[2].userId[2]) {
		t.Fatalf("expected user id %d, got %d", ticketData[2].userId[2], tickets[12].MatchmakingUserID)
	}
}

func Test_GetMatchmakingTickets_FilterUser_TicketsReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	tickets, err := q.GetMatchmakingTickets(context.Background(), GetMatchmakingTicketsParams{
		MatchmakingUser: MatchmakingUserParams{
			ID: sql.NullInt64{
				Int64: ticketData[0].userId[0],
				Valid: true,
			},
		},
		Limit:      10,
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking tickets: %v", err)
	}
	if len(tickets) != 4 {
		t.Fatalf("expected 4 tickets, got %d", len(tickets))
	}
	if tickets[0].TicketID != uint64(ticketIds[0]) {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], tickets[0].TicketID)
	}
	if tickets[0].MatchmakingUserID != uint64(ticketData[0].userId[0]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], tickets[0].MatchmakingUserID)
	}
	if tickets[0].ArenaID != uint64(ticketData[0].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[0], tickets[0].ArenaID)
	}
	if tickets[1].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], tickets[1].ArenaID)
	}
	if tickets[2].MatchmakingUserID != uint64(ticketData[0].userId[1]) {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], tickets[2].MatchmakingUserID)
	}
}

func Test_GetMatchmakingTickets_FilterArena_TicketsReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	tickets, err := q.GetMatchmakingTickets(context.Background(), GetMatchmakingTicketsParams{
		Arena: ArenaParams{
			ID: sql.NullInt64{
				Int64: ticketData[1].arenaId[1],
				Valid: true,
			},
		},
		Limit:      10,
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking tickets: %v", err)
	}
	if len(tickets) != 6 {
		t.Fatalf("expected 6 tickets, got %d", len(tickets))
	}
	if got, want := tickets[0].TicketID, uint64(ticketIds[1]); got != want {
		t.Fatalf("expected ticket id %d, got %d", want, got)
	}
	if got, want := tickets[0].MatchmakingUserID, uint64(ticketData[1].userId[0]); got != want {
		t.Fatalf("expected user id %d, got %d", want, got)
	}
	if got, want := tickets[0].ArenaID, uint64(ticketData[1].arenaId[0]); got != want {
		t.Fatalf("expected arena id %d, got %d", want, got)
	}
	if got, want := tickets[1].TicketID, uint64(ticketIds[1]); got != want {
		t.Fatalf("expected ticket id %d, got %d", want, got)
	}
	if got, want := tickets[1].MatchmakingUserID, uint64(ticketData[1].userId[0]); got != want {
		t.Fatalf("expected user id %d, got %d", want, got)
	}
	if got, want := tickets[1].ArenaID, uint64(ticketData[1].arenaId[1]); got != want {
		t.Fatalf("expected arena id %d, got %d", want, got)
	}
	if got, want := tickets[2].TicketID, uint64(ticketIds[1]); got != want {
		t.Fatalf("expected ticket id %d, got %d", want, got)
	}
	if got, want := tickets[2].MatchmakingUserID, uint64(ticketData[1].userId[0]); got != want {
		t.Fatalf("expected user id %d, got %d", want, got)
	}
	if got, want := tickets[2].ArenaID, uint64(ticketData[1].arenaId[2]); got != want {
		t.Fatalf("expected arena id %d, got %d", want, got)
	}
}

func Test_GetMatchmakingTickets_FilterEndedStatus_NoTicketsReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, _, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	tickets, err := q.GetMatchmakingTickets(context.Background(), GetMatchmakingTicketsParams{
		Statuses:   []string{"ENDED"},
		Limit:      10,
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking tickets: %v", err)
	}
	if len(tickets) != 0 {
		t.Fatalf("expected 0 tickets, got %d", len(tickets))
	}
}

func Test_GetMatchmakingTickets_FilterPendingStatus_TicketsReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, _, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	tickets, err := q.GetMatchmakingTickets(context.Background(), GetMatchmakingTicketsParams{
		Statuses:   []string{"PENDING"},
		Limit:      10,
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking tickets: %v", err)
	}
	if len(tickets) != 13 {
		t.Fatalf("expected 13 tickets, got %d", len(tickets))
	}
}

func Test_GetMatchmakingTickets_FilterMatchmakingUserAndArenaNoIntersection_NoTicketsReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	tickets, err := q.GetMatchmakingTickets(context.Background(), GetMatchmakingTicketsParams{
		MatchmakingUser: MatchmakingUserParams{
			ID: sql.NullInt64{
				Int64: ticketData[0].userId[0],
				Valid: true,
			},
		},
		Arena: ArenaParams{
			ID: sql.NullInt64{
				Int64: ticketData[1].arenaId[1],
				Valid: true,
			},
		},
		Limit:      10,
		UserLimit:  10,
		ArenaLimit: 10,
	})
	if err != nil {
		t.Fatalf("could not get matchmaking tickets: %v", err)
	}
	if len(tickets) != 0 {
		t.Fatalf("expected 0 tickets, got %d", len(tickets))
	}
}

// TODO: UpdateMatchmakingTicket tests go here

func Test_GetMatch_ByID_MatchExists(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	matchRows, err := q.GetMatch(context.Background(), GetMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get match: %v", err)
	}
	if len(matchRows) != 10 {
		t.Fatalf("expected 10 tickets, got %d", len(matchRows))
	}
	if matchRows[0].MatchID != uint64(matchId) {
		t.Fatalf("expected match id %d, got %d", matchId, matchRows[0].MatchID)
	}
	if matchRows[0].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matchRows[0].ArenaID)
	}
	if matchRows[0].TicketID.Int64 != ticketIds[0] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matchRows[0].TicketID.Int64)
	}
	if matchRows[0].TicketStatus.String != "MATCHED" {
		t.Fatalf("expected ticket status MATCHED, got %s", matchRows[0].TicketStatus.String)
	}
	if matchRows[0].MatchmakingUserID.Int64 != ticketData[0].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matchRows[0].MatchmakingUserID.Int64)
	}
	if matchRows[0].TicketArenaID.Int64 != ticketData[0].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[0], matchRows[0].TicketArenaID.Int64)
	}
	if matchRows[1].TicketArenaID.Int64 != ticketData[0].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[1], matchRows[1].TicketArenaID.Int64)
	}
	if matchRows[2].MatchmakingUserID.Int64 != ticketData[0].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matchRows[2].MatchmakingUserID.Int64)
	}
	if matchRows[4].TicketID.Int64 != ticketIds[1] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[1], matchRows[4].TicketID.Int64)
	}
	if matchRows[4].MatchmakingUserID.Int64 != ticketData[1].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[0], matchRows[4].MatchmakingUserID.Int64)
	}
	if matchRows[4].TicketArenaID.Int64 != ticketData[1].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[0], matchRows[4].TicketArenaID.Int64)
	}
	if matchRows[5].TicketArenaID.Int64 != ticketData[1].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[1], matchRows[5].TicketArenaID.Int64)
	}
	if matchRows[6].TicketArenaID.Int64 != ticketData[1].arenaId[2] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[2], matchRows[6].TicketArenaID.Int64)
	}
	if matchRows[7].MatchmakingUserID.Int64 != ticketData[1].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[1], matchRows[7].MatchmakingUserID.Int64)
	}
}

func Test_GetMatches_NoFilters_MatchesReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId1, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId1, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[2].arenaId[0], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId2, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of last ticket
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id = ?", matchId2, ticketIds[2])
	if err != nil {
		t.Fatalf("could not update ticket: %v", err)
	}
	matches, err := q.GetMatches(context.Background(), GetMatchesParams{
		Limit:       10,
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get matches: %v", err)
	}
	if len(matches) != 13 {
		t.Fatalf("expected 13 matches, got %d", len(matches))
	}
	if matches[0].MatchID != uint64(matchId1) {
		t.Fatalf("expected match id %d, got %d", matchId1, matches[0].MatchID)
	}
	if matches[0].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matches[0].ArenaID)
	}
	if matches[0].TicketID.Int64 != ticketIds[0] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matches[0].TicketID.Int64)
	}
	if matches[0].TicketStatus.String != "MATCHED" {
		t.Fatalf("expected ticket status MATCHED, got %s", matches[0].TicketStatus.String)
	}
	if matches[0].MatchmakingUserID.Int64 != ticketData[0].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matches[0].MatchmakingUserID.Int64)
	}
	if matches[0].TicketArenaID.Int64 != ticketData[0].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[0], matches[0].TicketArenaID.Int64)
	}
	if matches[1].TicketArenaID.Int64 != ticketData[0].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[1], matches[1].TicketArenaID.Int64)
	}
	if matches[2].MatchmakingUserID.Int64 != ticketData[0].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matches[2].MatchmakingUserID.Int64)
	}
	if matches[4].TicketID.Int64 != ticketIds[1] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[1], matches[4].TicketID.Int64)
	}
	if matches[4].MatchmakingUserID.Int64 != ticketData[1].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[0], matches[4].MatchmakingUserID.Int64)
	}
	if matches[4].TicketArenaID.Int64 != ticketData[1].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[0], matches[4].TicketArenaID.Int64)
	}
	if matches[5].TicketArenaID.Int64 != ticketData[1].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[1], matches[5].TicketArenaID.Int64)
	}
	if matches[6].TicketArenaID.Int64 != ticketData[1].arenaId[2] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[2], matches[6].TicketArenaID.Int64)
	}
	if matches[7].MatchmakingUserID.Int64 != ticketData[1].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[1], matches[7].MatchmakingUserID.Int64)
	}
	if matches[10].MatchID != uint64(matchId2) {
		t.Fatalf("expected match id %d, got %d", matchId2, matches[10].MatchID)
	}
	if matches[10].ArenaID != uint64(ticketData[2].arenaId[0]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[2].arenaId[0], matches[10].ArenaID)
	}
	if matches[10].TicketID.Int64 != ticketIds[2] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[2], matches[10].TicketID.Int64)
	}
	if matches[10].TicketStatus.String != "MATCHED" {
		t.Fatalf("expected ticket status MATCHED, got %s", matches[10].TicketStatus.String)
	}
	if matches[10].MatchmakingUserID.Int64 != ticketData[2].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[2].userId[0], matches[10].MatchmakingUserID.Int64)
	}
	if matches[10].TicketArenaID.Int64 != ticketData[2].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[2].arenaId[0], matches[10].TicketArenaID.Int64)
	}
	if matches[11].MatchmakingUserID.Int64 != ticketData[2].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[2].userId[1], matches[11].MatchmakingUserID.Int64)
	}
	if matches[12].MatchmakingUserID.Int64 != ticketData[2].userId[2] {
		t.Fatalf("expected user id %d, got %d", ticketData[2].userId[2], matches[12].MatchmakingUserID.Int64)
	}
}

func Test_GetMatches_FilterArena_MatchesReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId1, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId1, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[2].arenaId[0], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId2, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of last ticket
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id = ?", matchId2, ticketIds[2])
	if err != nil {
		t.Fatalf("could not update ticket: %v", err)
	}
	matches, err := q.GetMatches(context.Background(), GetMatchesParams{
		Arena: ArenaParams{
			ID: sql.NullInt64{
				Int64: ticketData[0].arenaId[1],
				Valid: true,
			},
		},
		Limit:       10,
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get matches: %v", err)
	}
	if len(matches) != 10 {
		t.Fatalf("expected 10 matches, got %d", len(matches))
	}
	if matches[0].MatchID != uint64(matchId1) {
		t.Fatalf("expected match id %d, got %d", matchId1, matches[0].MatchID)
	}
	if matches[0].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matches[0].ArenaID)
	}
	if matches[0].TicketID.Int64 != ticketIds[0] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matches[0].TicketID.Int64)
	}
	if matches[0].TicketStatus.String != "MATCHED" {
		t.Fatalf("expected ticket status MATCHED, got %s", matches[0].TicketStatus.String)
	}
	if matches[0].MatchmakingUserID.Int64 != ticketData[0].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matches[0].MatchmakingUserID.Int64)
	}
	if matches[0].TicketArenaID.Int64 != ticketData[0].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[0], matches[0].TicketArenaID.Int64)
	}
	if matches[1].TicketArenaID.Int64 != ticketData[0].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[1], matches[1].TicketArenaID.Int64)
	}
	if matches[2].MatchmakingUserID.Int64 != ticketData[0].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matches[2].MatchmakingUserID.Int64)
	}
	if matches[4].TicketID.Int64 != ticketIds[1] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[1], matches[4].TicketID.Int64)
	}
	if matches[4].MatchmakingUserID.Int64 != ticketData[1].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[0], matches[4].MatchmakingUserID.Int64)
	}
	if matches[4].TicketArenaID.Int64 != ticketData[1].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[0], matches[4].TicketArenaID.Int64)
	}
	if matches[5].TicketArenaID.Int64 != ticketData[1].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[1], matches[5].TicketArenaID.Int64)
	}
	if matches[6].TicketArenaID.Int64 != ticketData[1].arenaId[2] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[2], matches[6].TicketArenaID.Int64)
	}
	if matches[7].MatchmakingUserID.Int64 != ticketData[1].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[1], matches[7].MatchmakingUserID.Int64)
	}
}

func Test_GetMatches_FilterUser_MatchesReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId1, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId1, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[2].arenaId[0], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId2, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of last ticket
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id = ?", matchId2, ticketIds[2])
	if err != nil {
		t.Fatalf("could not update ticket: %v", err)
	}
	matches, err := q.GetMatches(context.Background(), GetMatchesParams{
		MatchmakingUser: MatchmakingUserParams{
			ID: sql.NullInt64{
				Int64: ticketData[0].userId[0],
				Valid: true,
			},
		},
		Limit:       10,
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get matches: %v", err)
	}
	if len(matches) != 10 {
		t.Fatalf("expected 10 matches, got %d", len(matches))
	}
	if matches[0].MatchID != uint64(matchId1) {
		t.Fatalf("expected match id %d, got %d", matchId1, matches[0].MatchID)
	}
	if matches[0].ArenaID != uint64(ticketData[0].arenaId[1]) {
		t.Fatalf("expected arena id %d, got %d", ticketData[0].arenaId[1], matches[0].ArenaID)
	}
	if matches[0].TicketID.Int64 != ticketIds[0] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[0], matches[0].TicketID.Int64)
	}
	if matches[0].TicketStatus.String != "MATCHED" {
		t.Fatalf("expected ticket status MATCHED, got %s", matches[0].TicketStatus.String)
	}
	if matches[0].MatchmakingUserID.Int64 != ticketData[0].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[0], matches[0].MatchmakingUserID.Int64)
	}
	if matches[0].TicketArenaID.Int64 != ticketData[0].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[0], matches[0].TicketArenaID.Int64)
	}
	if matches[1].TicketArenaID.Int64 != ticketData[0].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[0].arenaId[1], matches[1].TicketArenaID.Int64)
	}
	if matches[2].MatchmakingUserID.Int64 != ticketData[0].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[0].userId[1], matches[2].MatchmakingUserID.Int64)
	}
	if matches[4].TicketID.Int64 != ticketIds[1] {
		t.Fatalf("expected ticket id %d, got %d", ticketIds[1], matches[4].TicketID.Int64)
	}
	if matches[4].MatchmakingUserID.Int64 != ticketData[1].userId[0] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[0], matches[4].MatchmakingUserID.Int64)
	}
	if matches[4].TicketArenaID.Int64 != ticketData[1].arenaId[0] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[0], matches[4].TicketArenaID.Int64)
	}
	if matches[5].TicketArenaID.Int64 != ticketData[1].arenaId[1] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[1], matches[5].TicketArenaID.Int64)
	}
	if matches[6].TicketArenaID.Int64 != ticketData[1].arenaId[2] {
		t.Fatalf("expected ticket arena id %d, got %d", ticketData[1].arenaId[2], matches[6].TicketArenaID.Int64)
	}
	if matches[7].MatchmakingUserID.Int64 != ticketData[1].userId[1] {
		t.Fatalf("expected user id %d, got %d", ticketData[1].userId[1], matches[7].MatchmakingUserID.Int64)
	}
}

func Test_GetMatches_FilterMatchmakingUserAndArenaThatDontIntersect_NoMatchesReturned(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId1, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId1, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[2].arenaId[0], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId2, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of last ticket
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id = ?", matchId2, ticketIds[2])
	if err != nil {
		t.Fatalf("could not update ticket: %v", err)
	}
	matches, err := q.GetMatches(context.Background(), GetMatchesParams{
		MatchmakingUser: MatchmakingUserParams{
			ID: sql.NullInt64{
				Int64: ticketData[1].userId[0],
				Valid: true,
			},
		},
		Arena: ArenaParams{
			ID: sql.NullInt64{
				Int64: ticketData[0].arenaId[0],
				Valid: true,
			},
		},
		Limit:       10,
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get matches: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func Test_StartMatch_ByIDValidStartTime_MatchUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.StartMatch(context.Background(), StartMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		LockTime:  time.Now().Add(-time.Hour),
		StartTime: time.Now().Add(-time.Hour),
	})
	if err != nil {
		t.Fatalf("could not start match: %v", err)
	}
	match, err := q.GetMatch(context.Background(), GetMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get match: %v", err)
	}
	if match[0].StartedAt.Time.IsZero() {
		t.Fatalf("expected non-zero start time, got zero")
	}
	if match[0].MatchStatus != "STARTED" {
		t.Fatalf("expected match status STARTED, got %s", match[0].MatchStatus)
	}
}

func Test_StartMatch_ByIDStartTimeAlreadySet_MatchNotUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data, locked_at, started_at) VALUES (?, ?, ?, ?)", ticketData[0].arenaId[1], "{}", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.StartMatch(context.Background(), StartMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		LockTime:  time.Now(),
		StartTime: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not start match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_StartMatch_ByTicketIDValidStartTime_MatchUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.StartMatch(context.Background(), StartMatchParams{
		Match: MatchParams{
			MatchmakingTicket: MatchmakingTicketParams{
				ID: conversion.Int64ToSqlNullInt64(&ticketIds[0]),
			},
		},
		LockTime:  time.Now().Add(-time.Hour),
		StartTime: time.Now().Add(-time.Hour),
	})
	if err != nil {
		t.Fatalf("could not start match: %v", err)
	}
	match, err := q.GetMatch(context.Background(), GetMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get match: %v", err)
	}
	if match[0].StartedAt.Time.IsZero() {
		t.Fatalf("expected non-zero start time, got zero")
	}
	if match[0].MatchStatus != "STARTED" {
		t.Fatalf("expected match status STARTED, got %s", match[0].MatchStatus)
	}
}

func Test_StartMatch_ByMatchmakingUserIDValidStartTime_MatchUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.StartMatch(context.Background(), StartMatchParams{
		Match: MatchParams{
			MatchmakingTicket: MatchmakingTicketParams{
				MatchmakingUser: MatchmakingUserParams{
					ID: conversion.Int64ToSqlNullInt64(&ticketData[0].userId[0]),
				},
			},
		},
		LockTime:  time.Now().Add(-time.Hour),
		StartTime: time.Now().Add(-time.Hour),
	})
	if err != nil {
		t.Fatalf("could not start match: %v", err)
	}
	match, err := q.GetMatch(context.Background(), GetMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get match: %v", err)
	}
	if match[0].StartedAt.Time.IsZero() {
		t.Fatalf("expected non-zero start time, got zero")
	}
	if match[0].MatchStatus != "STARTED" {
		t.Fatalf("expected match status STARTED, got %s", match[0].MatchStatus)
	}
}

func Test_StartMatch_ByTicketIDMatchDoesntExist_NoRowsAffected(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.StartMatch(context.Background(), StartMatchParams{
		Match: MatchParams{
			MatchmakingTicket: MatchmakingTicketParams{
				ID: conversion.Int64ToSqlNullInt64(conversion.ValueToPointer(int64(999999999))),
			},
		},
		LockTime:  time.Now(),
		StartTime: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not start match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_StartMatch_ByMatchmakingUserIDMatchDoesntExist_NoRowsAffected(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.StartMatch(context.Background(), StartMatchParams{
		Match: MatchParams{
			MatchmakingTicket: MatchmakingTicketParams{
				MatchmakingUser: MatchmakingUserParams{
					ID: conversion.Int64ToSqlNullInt64(conversion.ValueToPointer(int64(999999999))),
				},
			},
		},
		LockTime:  time.Now(),
		StartTime: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not start match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_EndMatch_ByIDValidEndTime_MatchUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data, locked_at, started_at) VALUES (?, ?, ?, ?)", ticketData[0].arenaId[1], "{}", time.Now().Add(-time.Hour), time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.EndMatch(context.Background(), EndMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		EndTime: time.Now().Add(-time.Minute),
	})
	if err != nil {
		t.Fatalf("could not end match: %v", err)
	}
	match, err := q.GetMatch(context.Background(), GetMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		TicketLimit: 10,
		UserLimit:   10,
		ArenaLimit:  10,
	})
	if err != nil {
		t.Fatalf("could not get match: %v", err)
	}
	if match[0].EndedAt.Time.IsZero() {
		t.Fatalf("expected non-zero end time, got zero")
	}
	if match[0].MatchStatus != "ENDED" {
		t.Fatalf("expected match status ENDED, got %s", match[0].MatchStatus)
	}
}

func Test_EndMatch_ByIDEndedAtAlreadySet_MatchNotUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data, locked_at, started_at, ended_at) VALUES (?, ?, ?, ?, ?)", ticketData[0].arenaId[1], "{}", time.Now(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.EndMatch(context.Background(), EndMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		EndTime: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not end match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_EndMatch_ByIDStartTimeNotSet_MatchNotUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data, locked_at) VALUES (?, ?, ?)", ticketData[0].arenaId[1], "{}", time.Now())
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.EndMatch(context.Background(), EndMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		EndTime: time.Now(),
	})
	if err != nil {
		t.Fatalf("could not end match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_EndMatch_ByIDEndTimeBeforeStartTime_MatchNotUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	ticketIds, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	startTime := time.Now()
	endTime := startTime.Add(-time.Second)
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data, locked_at, started_at) VALUES (?, ?, ?, ?)", ticketData[0].arenaId[1], "{}", time.Now(), startTime)
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	// Set match id of first two tickets
	_, err = q.db.ExecContext(context.Background(), "UPDATE matchmaking_ticket SET matchmaking_match_id = ? WHERE id IN (?, ?)", matchId, ticketIds[0], ticketIds[1])
	if err != nil {
		t.Fatalf("could not update tickets: %v", err)
	}
	result, err = q.EndMatch(context.Background(), EndMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		EndTime: endTime,
	})
	if err != nil {
		t.Fatalf("could not end match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateMatch_ByIDValidData_MatchUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	data := map[string]interface{}{
		"key": "value",
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("could not marshal data: %v", err)
	}
	result, err = q.UpdateMatch(context.Background(), UpdateMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		Data: dataBytes,
	})
	if err != nil {
		t.Fatalf("could not update match: %v", err)
	}
	var newDataBytes json.RawMessage
	err = q.db.QueryRowContext(context.Background(), "SELECT data FROM matchmaking_match WHERE id = ?", matchId).Scan(&newDataBytes)
	if err != nil {
		t.Fatalf("could not get match data: %v", err)
	}
	var newData map[string]interface{}
	err = json.Unmarshal(newDataBytes, &newData)
	if err != nil {
		t.Fatalf("could not unmarshal data: %v", err)
	}
	if newData["key"] != "value" {
		t.Fatalf("expected key value, got %v", newData["key"])
	}
}

func Test_UpdateMatch_ByIDMatchDoesntExist_NoRowsAffected(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	data := map[string]interface{}{
		"key": "value",
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("could not marshal data: %v", err)
	}
	result, err := q.UpdateMatch(context.Background(), UpdateMatchParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(conversion.ValueToPointer(int64(999999999))),
		},
		Data: dataBytes,
	})
	if err != nil {
		t.Fatalf("could not update match: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_SetMatchPrivateServer_ByIDValidPrivateServerID_MatchUpdated(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	privateServerID := "41.41.41.41"
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data) VALUES (?, ?)", ticketData[0].arenaId[1], "{}")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.SetMatchPrivateServer(context.Background(), SetMatchPrivateServerParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		PrivateServerID: privateServerID,
	})
	if err != nil {
		t.Fatalf("could not set match private server: %v", err)
	}
	var actualPrivateServerID string
	err = q.db.QueryRowContext(context.Background(), "SELECT private_server_id FROM matchmaking_match WHERE id = ?", matchId).Scan(&actualPrivateServerID)
	if err != nil {
		t.Fatalf("could not get match private server id: %v", err)
	}
	if actualPrivateServerID != privateServerID {
		t.Fatalf("expected private server id %s, got %s", privateServerID, actualPrivateServerID)
	}
}

func Test_SetMatchPrivateServer_ByIDPrivateServerAlreadySet_NoRowsAffected(t *testing.T) {
	tx := server.Connect(t)
	q := New(tx)
	_, ticketData, err := createTestTickets(q)
	if err != nil {
		t.Fatalf("could not create test tickets: %v", err)
	}
	result, err := q.db.ExecContext(context.Background(), "INSERT INTO matchmaking_match (matchmaking_arena_id, data, private_server_id) VALUES (?, ?, ?)", ticketData[0].arenaId[1], "{}", "41.41.41.41")
	if err != nil {
		t.Fatalf("could not create match: %v", err)
	}
	matchId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	result, err = q.SetMatchPrivateServer(context.Background(), SetMatchPrivateServerParams{
		Match: MatchParams{
			ID: conversion.Int64ToSqlNullInt64(&matchId),
		},
		PrivateServerID: "281.281.281.281",
	})
	if err != nil {
		t.Fatalf("could not set match private server: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}
