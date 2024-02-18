package tournament

import (
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
)

func TestGetTournamentStartDateDaily(t *testing.T) {
	s := Service{}
	currentTime := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_DAILY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateDailyAtNoon(t *testing.T) {
	s := Service{
		dailyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_DAILY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateDailyAtNoonTheNextDay(t *testing.T) {
	s := Service{
		dailyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-02T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_DAILY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeekly(t *testing.T) {
	s := Service{
		weeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2019-12-30T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeeklyAtMonday(t *testing.T) {
	s := Service{
		weeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2023, 10, 9, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-09T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeeklyAtSunday(t *testing.T) {
	s := Service{
		weeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2023, 10, 8, 23, 59, 59, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-02T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeeklyAtMondayAtNoon(t *testing.T) {
	s := Service{
		weeklyTournamentDay:    time.Monday,
		weeklyTournamentMinute: 720,
	}
	currentTime := time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-09T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateWeeklyAtSundayAtNoonCurrentlyBeforeNoon(t *testing.T) {
	s := Service{
		weeklyTournamentDay:    time.Sunday,
		weeklyTournamentMinute: 720,
	}
	currentTime := time.Date(2023, 10, 8, 10, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_WEEKLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateMonthly(t *testing.T) {
	s := Service{
		monthlyTournamentDay: 1,
	}
	currentTime := time.Date(2020, 1, 5, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_MONTHLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateMonthlyAtDayTenAtNoon(t *testing.T) {
	s := Service{
		monthlyTournamentDay:    10,
		monthlyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-10T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_MONTHLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateMonthlyAtNoon(t *testing.T) {
	s := Service{
		monthlyTournamentDay:    1,
		monthlyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 5, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_MONTHLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func TestGetTournamentStartDateMonthlyAtNoonCurrentlyBeforeNoon(t *testing.T) {
	s := Service{
		monthlyTournamentDay:    1,
		monthlyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2019-12-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := s.GetTournamentStartDate(currentTime, api.TournamentInterval_MONTHLY)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

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
