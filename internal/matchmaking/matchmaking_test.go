package matchmaking

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

func Test_unmarshalArena_ValidArena_ValidApiArena(t *testing.T) {
	arena := model.MatchmakingArena{
		ID:                  1,
		Name:                "Arena 1",
		MinPlayers:          2,
		MaxPlayersPerTicket: 2,
		MaxPlayers:          4,
		Data:                json.RawMessage(`{"key": "value"}`),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	apiArena, err := unmarshalArena(arena)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	data, err := conversion.RawJsonToProtobufStruct(arena.Data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if got, want := apiArena.Id, arena.ID; got != want {
		t.Errorf("Expected ID %d, got %d", want, got)
	}
	if got, want := apiArena.Name, arena.Name; got != want {
		t.Errorf("Expected Name %s, got %s", want, got)
	}
	if got, want := apiArena.MinPlayers, arena.MinPlayers; got != want {
		t.Errorf("Expected MinPlayers %d, got %d", want, got)
	}
	if got, want := apiArena.MaxPlayersPerTicket, arena.MaxPlayersPerTicket; got != want {
		t.Errorf("Expected MaxPlayersPerTicket %d, got %d", want, got)
	}
	if got, want := apiArena.MaxPlayers, arena.MaxPlayers; got != want {
		t.Errorf("Expected MaxPlayers %d, got %d", want, got)
	}
	if got, want := apiArena.Data, data; !reflect.DeepEqual(got, want) {
		t.Errorf("Expected Data to be %v, got %v", want, got)
	}
	if got, want := apiArena.CreatedAt.AsTime(), arena.CreatedAt; !got.Equal(want) {
		t.Errorf("Expected CreatedAt %v, got %v", want, got)
	}
	if got, want := apiArena.UpdatedAt.AsTime(), arena.UpdatedAt; !got.Equal(want) {
		t.Errorf("Expected UpdatedAt %v, got %v", want, got)
	}
}

func Test_unmarshalMatchmakingUser_ValidMatchmakingUser_ValidApiMatchmakingUser(t *testing.T) {
	mu := model.MatchmakingUser{
		ID:           1,
		ClientUserID: 102,
		Elo:          1500,
		Data:         json.RawMessage(`{"key": "value"}`),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	apiMu, err := unmarshalMatchmakingUser(mu)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	data, err := conversion.RawJsonToProtobufStruct(mu.Data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if got, want := apiMu.Id, mu.ID; got != want {
		t.Errorf("Expected ID %d, got %d", want, got)
	}
	if got, want := apiMu.ClientUserId, mu.ClientUserID; got != want {
		t.Errorf("Expected ClientUserId %d, got %d", want, got)
	}
	if got, want := apiMu.Elo, mu.Elo; got != want {
		t.Errorf("Expected Elo %d, got %d", want, got)
	}
	if got, want := apiMu.Data, data; !reflect.DeepEqual(got, want) {
		t.Errorf("Expected Data to be %v, got %v", want, got)
	}
	if got, want := apiMu.CreatedAt.AsTime(), mu.CreatedAt; !got.Equal(want) {
		t.Errorf("Expected CreatedAt %v, got %v", want, got)
	}
	if got, want := apiMu.UpdatedAt.AsTime(), mu.UpdatedAt; !got.Equal(want) {
		t.Errorf("Expected UpdatedAt %v, got %v", want, got)
	}
}

func Test_unmarshalMatchmakingTicket_ValidMatchmakingTicket_ValidApiMatchmakingTicket(t *testing.T) {
	mt := []model.MatchmakingTicketWithUserAndArena{
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        4,
			ClientUserID:             102,
			Elo:                      1500,
			UserNumber:               1,
			UserData:                 json.RawMessage(`{"userKey":"value1"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  5,
			ArenaName:                "Arena 1",
			ArenaMinPlayers:          2,
			ArenaMaxPlayersPerTicket: 2,
			ArenaMaxPlayers:          4,
			ArenaNumber:              1,
			ArenaData:                json.RawMessage(`{"arenaKey":"value1"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        4,
			ClientUserID:             102,
			Elo:                      1500,
			UserNumber:               1,
			UserData:                 json.RawMessage(`{"userKey":"value1"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  6,
			ArenaName:                "Arena 2",
			ArenaMinPlayers:          3,
			ArenaMaxPlayersPerTicket: 3,
			ArenaMaxPlayers:          6,
			ArenaNumber:              2,
			ArenaData:                json.RawMessage(`{"arenaKey":"value2"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        5,
			ClientUserID:             103,
			Elo:                      1600,
			UserNumber:               2,
			UserData:                 json.RawMessage(`{"userKey":"value2"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  5,
			ArenaName:                "Arena 1",
			ArenaMinPlayers:          2,
			ArenaMaxPlayersPerTicket: 2,
			ArenaMaxPlayers:          4,
			ArenaNumber:              1,
			ArenaData:                json.RawMessage(`{"arenaKey":"value1"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        5,
			ClientUserID:             103,
			Elo:                      1600,
			UserNumber:               2,
			UserData:                 json.RawMessage(`{"userKey":"value2"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  6,
			ArenaName:                "Arena 2",
			ArenaMinPlayers:          3,
			ArenaMaxPlayersPerTicket: 3,
			ArenaMaxPlayers:          6,
			ArenaNumber:              2,
			ArenaData:                json.RawMessage(`{"arenaKey":"value2"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        6,
			ClientUserID:             104,
			Elo:                      1700,
			UserNumber:               3,
			UserData:                 json.RawMessage(`{"userKey":"value3"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  5,
			ArenaName:                "Arena 1",
			ArenaMinPlayers:          2,
			ArenaMaxPlayersPerTicket: 2,
			ArenaMaxPlayers:          4,
			ArenaNumber:              1,
			ArenaData:                json.RawMessage(`{"arenaKey":"value1"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        6,
			ClientUserID:             104,
			Elo:                      1700,
			UserNumber:               3,
			UserData:                 json.RawMessage(`{"userKey":"value3"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  6,
			ArenaName:                "Arena 2",
			ArenaMinPlayers:          3,
			ArenaMaxPlayersPerTicket: 3,
			ArenaMaxPlayers:          6,
			ArenaNumber:              2,
			ArenaData:                json.RawMessage(`{"arenaKey":"value2"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
	}
	apiTickets, err := unmarshalMatchmakingTicket(mt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ticketData, err := conversion.RawJsonToProtobufStruct(mt[0].TicketData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if got, want := apiTickets.Id, mt[0].TicketID; got != want {
		t.Errorf("Expected ID %d, got %d", want, got)
	}
	if got, want := *apiTickets.MatchId, uint64(mt[0].MatchmakingMatchID.Int64); got != want {
		t.Errorf("Expected MatchId %v, got %v", want, got)
	}
	if got, want := apiTickets.Status.String(), mt[0].Status; got != want {
		t.Errorf("Expected Status %s, got %s", want, got)
	}
	if got, want := apiTickets.Data, ticketData; !reflect.DeepEqual(got, want) {
		t.Errorf("Expected Data to be %v, got %v", want, got)
	}
	if got, want := apiTickets.ExpiresAt.AsTime(), mt[0].ExpiresAt; !got.Equal(want) {
		t.Errorf("Expected ExpiresAt %v, got %v", want, got)
	}
	if got, want := apiTickets.CreatedAt.AsTime(), mt[0].TicketCreatedAt; !got.Equal(want) {
		t.Errorf("Expected CreatedAt %v, got %v", want, got)
	}
	if got, want := apiTickets.UpdatedAt.AsTime(), mt[0].TicketUpdatedAt; !got.Equal(want) {
		t.Errorf("Expected UpdatedAt %v, got %v", want, got)
	}
	if got, want := len(apiTickets.MatchmakingUsers), 3; got != want {
		t.Fatalf("Expected %d MatchmakingUsers, got %d", want, got)
	}
	for i, user := range apiTickets.MatchmakingUsers {
		userIndex := i * 2
		userData, err := conversion.RawJsonToProtobufStruct(mt[userIndex].UserData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got, want := user.Id, mt[userIndex].MatchmakingUserID; got != want {
			t.Errorf("Expected User ID %d, got %d", want, got)
		}
		if got, want := user.ClientUserId, mt[userIndex].ClientUserID; got != want {
			t.Errorf("Expected ClientUserId %d, got %d", want, got)
		}
		if got, want := user.Elo, mt[userIndex].Elo; got != want {
			t.Errorf("Expected Elo %d, got %d", want, got)
		}
		if got, want := user.Data, userData; !reflect.DeepEqual(got, want) {
			t.Errorf("Expected User Data to be %v, got %v", want, got)
		}
		if got, want := user.CreatedAt.AsTime(), mt[userIndex].UserCreatedAt; !got.Equal(want) {
			t.Errorf("Expected User CreatedAt %v, got %v", want, got)
		}
		if got, want := user.UpdatedAt.AsTime(), mt[userIndex].UserUpdatedAt; !got.Equal(want) {
			t.Errorf("Expected User UpdatedAt %v, got %v", want, got)
		}
	}
	if got, want := len(apiTickets.Arenas), 2; got != want {
		t.Fatalf("Expected %d Arenas, got %d", want, got)
	}
	for i, arena := range apiTickets.Arenas {
		if got, want := arena.Id, mt[i].ArenaID; got != want {
			t.Errorf("Expected Arena ID %d, got %d", want, got)
		}
		if got, want := arena.Name, mt[i].ArenaName; got != want {
			t.Errorf("Expected Arena Name %s, got %s", want, got)
		}
		if got, want := arena.MinPlayers, mt[i].ArenaMinPlayers; got != want {
			t.Errorf("Expected Arena MinPlayers %d, got %d", want, got)
		}
		if got, want := arena.MaxPlayersPerTicket, mt[i].ArenaMaxPlayersPerTicket; got != want {
			t.Errorf("Expected Arena MaxPlayersPerTicket %d, got %d", want, got)
		}
		if got, want := arena.MaxPlayers, mt[i].ArenaMaxPlayers; got != want {
			t.Errorf("Expected Arena MaxPlayers %d, got %d", want, got)
		}
		arenaData, err := conversion.RawJsonToProtobufStruct(mt[i].ArenaData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got, want := arena.Data, arenaData; !reflect.DeepEqual(got, want) {
			t.Errorf("Expected Arena Data to be %v, got %v", want, got)
		}
		if got, want := arena.CreatedAt.AsTime(), mt[i].ArenaCreatedAt; !got.Equal(want) {
			t.Errorf("Expected Arena CreatedAt %v, got %v", want, got)
		}
		if got, want := arena.UpdatedAt.AsTime(), mt[i].ArenaUpdatedAt; !got.Equal(want) {
			t.Errorf("Expected Arena UpdatedAt %v, got %v", want, got)
		}
	}
}

func Test_unmarshalMatchmakingTickets_ValidMatchmakingTickets_ValidApiMatchmakingTickets(t *testing.T) {
	mt := []model.MatchmakingTicketWithUserAndArena{
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        4,
			ClientUserID:             102,
			Elo:                      1500,
			UserNumber:               1,
			UserData:                 json.RawMessage(`{"userKey":"value1"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  5,
			ArenaName:                "Arena 1",
			ArenaMinPlayers:          2,
			ArenaMaxPlayersPerTicket: 2,
			ArenaMaxPlayers:          4,
			ArenaNumber:              1,
			ArenaData:                json.RawMessage(`{"arenaKey":"value1"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        4,
			ClientUserID:             102,
			Elo:                      1500,
			UserNumber:               1,
			UserData:                 json.RawMessage(`{"userKey":"value1"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  6,
			ArenaName:                "Arena 2",
			ArenaMinPlayers:          3,
			ArenaMaxPlayersPerTicket: 3,
			ArenaMaxPlayers:          6,
			ArenaNumber:              2,
			ArenaData:                json.RawMessage(`{"arenaKey":"value2"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        5,
			ClientUserID:             103,
			Elo:                      1600,
			UserNumber:               2,
			UserData:                 json.RawMessage(`{"userKey":"value2"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  5,
			ArenaName:                "Arena 1",
			ArenaMinPlayers:          2,
			ArenaMaxPlayersPerTicket: 2,
			ArenaMaxPlayers:          4,
			ArenaNumber:              1,
			ArenaData:                json.RawMessage(`{"arenaKey":"value1"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        5,
			ClientUserID:             103,
			Elo:                      1600,
			UserNumber:               2,
			UserData:                 json.RawMessage(`{"userKey":"value2"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  6,
			ArenaName:                "Arena 2",
			ArenaMinPlayers:          3,
			ArenaMaxPlayersPerTicket: 3,
			ArenaMaxPlayers:          6,
			ArenaNumber:              2,
			ArenaData:                json.RawMessage(`{"arenaKey":"value2"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        6,
			ClientUserID:             104,
			Elo:                      1700,
			UserNumber:               3,
			UserData:                 json.RawMessage(`{"userKey":"value3"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  5,
			ArenaName:                "Arena 1",
			ArenaMinPlayers:          2,
			ArenaMaxPlayersPerTicket: 2,
			ArenaMaxPlayers:          4,
			ArenaNumber:              1,
			ArenaData:                json.RawMessage(`{"arenaKey":"value1"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 1,
			MatchmakingMatchID:       sql.NullInt64{Int64: 2, Valid: true},
			Status:                   "MATCHED",
			UserCount:                3,
			TicketData:               json.RawMessage(`{"ticketKey":"value"}`),
			ExpiresAt:                time.Now().Add(1 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        6,
			ClientUserID:             104,
			Elo:                      1700,
			UserNumber:               3,
			UserData:                 json.RawMessage(`{"userKey":"value3"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  6,
			ArenaName:                "Arena 2",
			ArenaMinPlayers:          3,
			ArenaMaxPlayersPerTicket: 3,
			ArenaMaxPlayers:          6,
			ArenaNumber:              2,
			ArenaData:                json.RawMessage(`{"arenaKey":"value2"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 2,
			MatchmakingMatchID:       sql.NullInt64{Int64: 0, Valid: false},
			Status:                   "MATCHED",
			UserCount:                2,
			TicketData:               json.RawMessage(`{"ticketKey":"value2"}`),
			ExpiresAt:                time.Now().Add(2 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        7,
			ClientUserID:             105,
			Elo:                      1800,
			UserNumber:               1,
			UserData:                 json.RawMessage(`{"userKey":"value4"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  7,
			ArenaName:                "Arena 3",
			ArenaMinPlayers:          4,
			ArenaMaxPlayersPerTicket: 4,
			ArenaMaxPlayers:          8,
			ArenaNumber:              3,
			ArenaData:                json.RawMessage(`{"arenaKey":"value3"}`),
			ArenaCreatedAt:           time.Now(),
			ArenaUpdatedAt:           time.Now(),
		},
		{
			TicketID:                 2,
			MatchmakingMatchID:       sql.NullInt64{Int64: 0, Valid: false},
			Status:                   "MATCHED",
			UserCount:                2,
			TicketData:               json.RawMessage(`{"ticketKey":"value2"}`),
			ExpiresAt:                time.Now().Add(2 * time.Hour),
			TicketCreatedAt:          time.Now(),
			TicketUpdatedAt:          time.Now(),
			MatchmakingUserID:        8,
			ClientUserID:             106,
			Elo:                      1900,
			UserNumber:               2,
			UserData:                 json.RawMessage(`{"userKey":"value5"}`),
			UserCreatedAt:            time.Now(),
			UserUpdatedAt:            time.Now(),
			ArenaID:                  7,
			ArenaName:                "Arena 3",
			ArenaMinPlayers:          4,
			ArenaMaxPlayersPerTicket: 4,
			ArenaMaxPlayers:          8,
		},
	}
	apiTickets, err := unmarshalMatchmakingTickets(mt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ticketData, err := conversion.RawJsonToProtobufStruct(mt[0].TicketData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if got, want := len(apiTickets), 2; got != want {
		t.Fatalf("Expected %d MatchmakingTickets, got %d", want, got)
	}
	// Ticket 1
	if got, want := apiTickets[0].Id, mt[0].TicketID; got != want {
		t.Errorf("Expected ID %d, got %d", want, got)
	}
	if got, want := *apiTickets[0].MatchId, uint64(mt[0].MatchmakingMatchID.Int64); got != want {
		t.Errorf("Expected MatchId %v, got %v", want, got)
	}
	if got, want := apiTickets[0].Status.String(), mt[0].Status; got != want {
		t.Errorf("Expected Status %s, got %s", want, got)
	}
	if got, want := apiTickets[0].Data, ticketData; !reflect.DeepEqual(got, want) {
		t.Errorf("Expected Data to be %v, got %v", want, got)
	}
	if got, want := apiTickets[0].ExpiresAt.AsTime(), mt[0].ExpiresAt; !got.Equal(want) {
		t.Errorf("Expected ExpiresAt %v, got %v", want, got)
	}
	if got, want := apiTickets[0].CreatedAt.AsTime(), mt[0].TicketCreatedAt; !got.Equal(want) {
		t.Errorf("Expected CreatedAt %v, got %v", want, got)
	}
	if got, want := apiTickets[0].UpdatedAt.AsTime(), mt[0].TicketUpdatedAt; !got.Equal(want) {
		t.Errorf("Expected UpdatedAt %v, got %v", want, got)
	}
	if got, want := len(apiTickets[0].MatchmakingUsers), 3; got != want {
		t.Fatalf("Expected %d MatchmakingUsers, got %d", want, got)
	}
	for i, user := range apiTickets[0].MatchmakingUsers {
		userIndex := i * 2 // Each user appears 2 times in the test data
		userData, err := conversion.RawJsonToProtobufStruct(mt[userIndex].UserData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got, want := user.Id, mt[userIndex].MatchmakingUserID; got != want {
			t.Errorf("Expected User ID %d, got %d", want, got)
		}
		if got, want := user.ClientUserId, mt[userIndex].ClientUserID; got != want {
			t.Errorf("Expected ClientUserId %d, got %d", want, got)
		}
		if got, want := user.Elo, mt[userIndex].Elo; got != want {
			t.Errorf("Expected Elo %d, got %d", want, got)
		}
		if got, want := user.Data, userData; !reflect.DeepEqual(got, want) {
			t.Errorf("Expected User Data to be %v, got %v", want, got)
		}
		if got, want := user.CreatedAt.AsTime(), mt[userIndex].UserCreatedAt; !got.Equal(want) {
			t.Errorf("Expected User CreatedAt %v, got %v", want, got)
		}
		if got, want := user.UpdatedAt.AsTime(), mt[userIndex].UserUpdatedAt; !got.Equal(want) {
			t.Errorf("Expected User UpdatedAt %v, got %v", want, got)
		}
	}
	if got, want := len(apiTickets[0].Arenas), 2; got != want {
		t.Fatalf("Expected %d Arenas, got %d", want, got)
	}
	for i, arena := range apiTickets[0].Arenas {
		if got, want := arena.Id, mt[i].ArenaID; got != want {
			t.Errorf("Expected Arena ID %d, got %d", want, got)
		}
		if got, want := arena.Name, mt[i].ArenaName; got != want {
			t.Errorf("Expected Arena Name %s, got %s", want, got)
		}
		if got, want := arena.MinPlayers, mt[i].ArenaMinPlayers; got != want {
			t.Errorf("Expected Arena MinPlayers %d, got %d", want, got)
		}
		if got, want := arena.MaxPlayersPerTicket, mt[i].ArenaMaxPlayersPerTicket; got != want {
			t.Errorf("Expected Arena MaxPlayersPerTicket %d, got %d", want, got)
		}
		if got, want := arena.MaxPlayers, mt[i].ArenaMaxPlayers; got != want {
			t.Errorf("Expected Arena MaxPlayers %d, got %d", want, got)
		}
		arenaData, err := conversion.RawJsonToProtobufStruct(mt[i].ArenaData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got, want := arena.Data, arenaData; !reflect.DeepEqual(got, want) {
			t.Errorf("Expected Arena Data to be %v, got %v", want, got)
		}
		if got, want := arena.CreatedAt.AsTime(), mt[i].ArenaCreatedAt; !got.Equal(want) {
			t.Errorf("Expected Arena CreatedAt %v, got %v", want, got)
		}
		if got, want := arena.UpdatedAt.AsTime(), mt[i].ArenaUpdatedAt; !got.Equal(want) {
			t.Errorf("Expected Arena UpdatedAt %v, got %v", want, got)
		}
	}
	// Ticket 2
	ticketData, err = conversion.RawJsonToProtobufStruct(mt[6].TicketData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if got, want := apiTickets[1].Id, mt[6].TicketID; got != want {
		t.Errorf("Expected ID %d, got %d", want, got)
	}
	if got, want := apiTickets[1].MatchId, (*uint64)(nil); got != want {
		t.Errorf("Expected MatchId %v, got %v", want, got)
	}
	if got, want := apiTickets[1].Status.String(), mt[6].Status; got != want {
		t.Errorf("Expected Status %s, got %s", want, got)
	}
	if got, want := apiTickets[1].Data, ticketData; !reflect.DeepEqual(got, want) {
		t.Errorf("Expected Data to be %v, got %v", want, got)
	}
	if got, want := apiTickets[1].ExpiresAt.AsTime(), mt[6].ExpiresAt; !got.Equal(want) {
		t.Errorf("Expected ExpiresAt %v, got %v", want, got)
	}
	if got, want := apiTickets[1].CreatedAt.AsTime(), mt[6].TicketCreatedAt; !got.Equal(want) {
		t.Errorf("Expected CreatedAt %v, got %v", want, got)
	}
	if got, want := apiTickets[1].UpdatedAt.AsTime(), mt[6].TicketUpdatedAt; !got.Equal(want) {
		t.Errorf("Expected UpdatedAt %v, got %v", want, got)
	}
	if got, want := len(apiTickets[1].MatchmakingUsers), 2; got != want {
		t.Fatalf("Expected %d MatchmakingUsers, got %d", want, got)
	}
	for i, user := range apiTickets[1].MatchmakingUsers {
		userData, err := conversion.RawJsonToProtobufStruct(mt[6+i].UserData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got, want := user.Id, mt[6+i].MatchmakingUserID; got != want {
			t.Errorf("Expected User ID %d, got %d", want, got)
		}
		if got, want := user.ClientUserId, mt[6+i].ClientUserID; got != want {
			t.Errorf("Expected ClientUserId %d, got %d", want, got)
		}
		if got, want := user.Elo, mt[6+i].Elo; got != want {
			t.Errorf("Expected Elo %d, got %d", want, got)
		}
		if got, want := user.Data, userData; !reflect.DeepEqual(got, want) {
			t.Errorf("Expected User Data to be %v, got %v", want, got)
		}
		if got, want := user.CreatedAt.AsTime(), mt[6+i].UserCreatedAt; !got.Equal(want) {
			t.Errorf("Expected User CreatedAt %v, got %v", want, got)
		}
		if got, want := user.UpdatedAt.AsTime(), mt[6+i].UserUpdatedAt; !got.Equal(want) {
			t.Errorf("Expected User UpdatedAt %v, got %v", want, got)
		}
	}
	if got, want := len(apiTickets[1].Arenas), 1; got != want {
		t.Fatalf("Expected %d Arenas, got %d", want, got)
	}
	for i, arena := range apiTickets[1].Arenas {
		arenaIndex := 6 + i // Each arena appears once in the test data for ticket 2
		if got, want := arena.Id, mt[arenaIndex].ArenaID; got != want {
			t.Errorf("Expected Arena ID %d, got %d", want, got)
		}
		if got, want := arena.Name, mt[arenaIndex].ArenaName; got != want {
			t.Errorf("Expected Arena Name %s, got %s", want, got)
		}
		if got, want := arena.MinPlayers, mt[arenaIndex].ArenaMinPlayers; got != want {
			t.Errorf("Expected Arena MinPlayers %d, got %d", want, got)
		}
		if got, want := arena.MaxPlayersPerTicket, mt[arenaIndex].ArenaMaxPlayersPerTicket; got != want {
			t.Errorf("Expected Arena MaxPlayersPerTicket %d, got %d", want, got)
		}
		if got, want := arena.MaxPlayers, mt[arenaIndex].ArenaMaxPlayers; got != want {
			t.Errorf("Expected Arena MaxPlayers %d, got %d", want, got)
		}
		arenaData, err := conversion.RawJsonToProtobufStruct(mt[arenaIndex].ArenaData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got, want := arena.Data, arenaData; !reflect.DeepEqual(got, want) {
			t.Errorf("Expected Arena Data to be %v, got %v", want, got)
		}
		if got, want := arena.CreatedAt.AsTime(), mt[arenaIndex].ArenaCreatedAt; !got.Equal(want) {
			t.Errorf("Expected Arena CreatedAt %v, got %v", want, got)
		}
		if got, want := arena.UpdatedAt.AsTime(), mt[arenaIndex].ArenaUpdatedAt; !got.Equal(want) {
			t.Errorf("Expected Arena UpdatedAt %v, got %v", want, got)
		}
	}
}

// TODO: UnmarshalMatch, UnmarshalMatches, etc.
