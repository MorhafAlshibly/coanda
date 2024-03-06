package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

func TestWipeTournaments(t *testing.T) {
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{
		"test": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	tournament := model.RankedTournament{
		ID:                  1,
		Name:                "test",
		UserID:              1,
		TournamentInterval:  model.TournamentTournamentIntervalDaily,
		Score:               1,
		Ranking:             1,
		Data:                raw,
		TournamentStartedAt: time.Now(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	tournamentMap, err := json.Marshal(tournament)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(string(tournamentMap))
}
