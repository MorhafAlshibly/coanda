package tournament

import (
	"time"

	"github.com/MorhafAlshibly/coanda/api"
)

type WipeTimes struct {
	DailyTournamentMinute   uint16
	WeeklyTournamentMinute  uint16
	WeeklyTournamentDay     time.Weekday
	MonthlyTournamentMinute uint16
	MonthlyTournamentDay    uint8
}

func GetStartTime(currentTime time.Time, interval api.TournamentInterval, times WipeTimes) time.Time {
	var startDate time.Time
	switch interval {
	case api.TournamentInterval_DAILY:
		startDate = currentTime.UTC().Truncate(time.Hour * 24).Add(time.Duration(times.DailyTournamentMinute) * time.Minute)
		if currentTime.UTC().Before(startDate) {
			startDate = startDate.Add(-24 * time.Hour)
		}
	case api.TournamentInterval_WEEKLY:
		startDate = currentTime.UTC().Truncate(time.Hour * 24).Add(time.Duration((int(times.WeeklyTournamentDay)-int(currentTime.UTC().Weekday())-7)%7) * 24 * time.Hour).Add(time.Duration(times.WeeklyTournamentMinute) * time.Minute)
		if currentTime.UTC().Before(startDate) {
			startDate = startDate.Add(-7 * 24 * time.Hour)
		}
	case api.TournamentInterval_MONTHLY:
		startDate = time.Date(currentTime.UTC().Year(), currentTime.UTC().Month(), int(times.MonthlyTournamentDay), 0, 0, 0, 0, time.UTC).Add(time.Duration(times.MonthlyTournamentMinute) * time.Minute)
		if currentTime.UTC().Before(startDate) {
			startDate = startDate.AddDate(0, -1, 0)
		}
	default:
		startDate = time.Unix(0, 0).UTC()
	}
	return startDate
}
