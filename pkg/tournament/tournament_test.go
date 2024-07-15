package tournament

import (
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
)

func Test_GetStartTime_StartDateDaily_StartOfDay(t *testing.T) {
	times := WipeTimes{}
	currentTime := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_DAILY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateDailyAtNoon_NoonOfDay(t *testing.T) {
	times := WipeTimes{
		DailyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_DAILY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateDailyAtNoonTheNextDay_NoonNextDay(t *testing.T) {
	times := WipeTimes{
		DailyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-02T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_DAILY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateWeekly_StartOfWeek(t *testing.T) {
	times := WipeTimes{
		WeeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2019-12-30T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_WEEKLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateWeeklyAtMonday_StartOfMonday(t *testing.T) {
	times := WipeTimes{
		WeeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2023, 10, 9, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-09T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_WEEKLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateWeeklyAtSunday_StartOfSunday(t *testing.T) {
	times := WipeTimes{
		WeeklyTournamentDay: time.Monday,
	}
	currentTime := time.Date(2023, 10, 8, 23, 59, 59, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-02T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_WEEKLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateWeeklyAtMondayAtNoon_MondayAtNoon(t *testing.T) {
	times := WipeTimes{
		WeeklyTournamentDay:    time.Monday,
		WeeklyTournamentMinute: 720,
	}
	currentTime := time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-09T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_WEEKLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateWeeklyAtSundayAtNoonCurrentlyBeforeNoon_SundayAtNoon(t *testing.T) {
	times := WipeTimes{
		WeeklyTournamentDay:    time.Sunday,
		WeeklyTournamentMinute: 720,
	}
	currentTime := time.Date(2023, 10, 8, 10, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_WEEKLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateMonthly_StartOfMonth(t *testing.T) {
	times := WipeTimes{
		MonthlyTournamentDay: 1,
	}
	currentTime := time.Date(2020, 1, 5, 1, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_MONTHLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateMonthlyAtDayTenAtNoon_DayTenAtNoon(t *testing.T) {
	times := WipeTimes{
		MonthlyTournamentDay:    10,
		MonthlyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-10T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_MONTHLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateMonthlyAtNoon_StartOfMonthAtNoon(t *testing.T) {
	times := WipeTimes{
		MonthlyTournamentDay:    1,
		MonthlyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 5, 12, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2020-01-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_MONTHLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}

func Test_GetStartTime_StartDateMonthlyAtNoonCurrentlyBeforeNoon_StartOfLastMonthAtNoon(t *testing.T) {
	times := WipeTimes{
		MonthlyTournamentDay:    1,
		MonthlyTournamentMinute: 720,
	}
	currentTime := time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)
	correctDate, err := time.Parse(time.RFC3339, "2019-12-01T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startDate := GetStartTime(currentTime, api.TournamentInterval_MONTHLY, times)
	if startDate != correctDate {
		t.Fatalf("Expected start date to be %s, got %s", correctDate, startDate)
	}
}
