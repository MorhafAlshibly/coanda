package tournament

import (
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
)

func TestGetTournamentStartDateDaily(t *testing.T) {
	s := Service{}
	currentTime := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	correctDate := "2020-01-01T00:00:00Z"
	startDate := s.getTournamentStartDate(currentTime, api.TournamentInterval_DAILY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeekly(t *testing.T) {
	s := Service{
		weeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	correctDate := "2019-12-30T00:00:00Z"
	startDate := s.getTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeeklyAtMonday(t *testing.T) {
	s := Service{
		weeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2023, 10, 9, 1, 0, 0, 0, time.UTC)
	correctDate := "2023-10-09T00:00:00Z"
	startDate := s.getTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeeklyAtSunday(t *testing.T) {
	s := Service{
		weeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2023, 10, 8, 23, 59, 59, 0, time.UTC)
	correctDate := "2023-10-02T00:00:00Z"
	startDate := s.getTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateMonthly(t *testing.T) {
	s := Service{
		monthlyTournamentDay: 1,
	}
	currentTime := time.Date(2020, 1, 5, 1, 0, 0, 0, time.UTC)
	correctDate := "2020-01-01T00:00:00Z"
	startDate := s.getTournamentStartDate(currentTime, api.TournamentInterval_MONTHLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}
