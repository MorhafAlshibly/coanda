package tournament

import (
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
)

func TestConvertTournamentIntervalUserIdToNullNameIntervalUserIDStartedAt(t *testing.T) {
	s := Service{}
	tournamentIntervalUserId := &api.TournamentIntervalUserId{
		Tournament: "test",
		Interval:   api.TournamentInterval_DAILY,
		UserId:     1,
	}
	nullNameIntervalUserIDStartedAt := s.convertTournamentIntervalUserIdToNullNameIntervalUserIDStartedAt(tournamentIntervalUserId)
	if nullNameIntervalUserIDStartedAt.Name != "test" {
		t.Fatalf("Expected name to be test, got %s", nullNameIntervalUserIDStartedAt.Name)
	}
	if nullNameIntervalUserIDStartedAt.TournamentInterval != "DAILY" {
		t.Fatalf("Expected tournament interval to be DAILY, got %s", nullNameIntervalUserIDStartedAt.TournamentInterval)
	}
	if nullNameIntervalUserIDStartedAt.UserID != 1 {
		t.Fatalf("Expected user ID to be 1, got %d", nullNameIntervalUserIDStartedAt.UserID)
	}
	if nullNameIntervalUserIDStartedAt.Valid != true {
		t.Fatal("Expected valid to be true")
	}
}
