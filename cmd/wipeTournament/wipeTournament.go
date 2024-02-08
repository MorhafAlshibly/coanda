package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	dsn                     = os.Getenv("DSN")
	dailyTournamentMinute   = os.Getenv("DAILY_TOURNAMENT_MINUTE")
	weeklyTournamentMinute  = os.Getenv("WEEKLY_TOURNAMENT_MINUTE")
	weeklyTournamentDay     = os.Getenv("WEEKLY_TOURNAMENT_DAY")
	monthlyTournamentMinute = os.Getenv("MONTHLY_TOURNAMENT_MINUTE")
	monthlyTournamentDay    = os.Getenv("MONTHLY_TOURNAMENT_DAY")
)

type WipeTournamentApp struct {
	tournamentService *tournament.Service
}

func NewWipeTournamentApp() *WipeTournamentApp {
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db := model.New(dbConn)
	return &WipeTournamentApp{
		tournamentService: tournament.NewService(
			tournament.WithSql(dbConn),
			tournament.WithDatabase(db),
			tournament.WithDailyTournamentMinute(uint16(stringToIntPanicOnError(dailyTournamentMinute))),
			tournament.WithWeeklyTournamentMinute(uint16(stringToIntPanicOnError(weeklyTournamentMinute))),
			tournament.WithWeeklyTournamentDay(time.Weekday(stringToIntPanicOnError(weeklyTournamentDay))),
			tournament.WithMonthlyTournamentMinute(uint16(stringToIntPanicOnError(monthlyTournamentMinute))),
			tournament.WithMonthlyTournamentDay(uint8(stringToIntPanicOnError(monthlyTournamentDay))),
		),
	}
}

func (a *WipeTournamentApp) handler(ctx context.Context) {
	a.WipeTournaments(ctx, api.TournamentInterval_DAILY)
	a.WipeTournaments(ctx, api.TournamentInterval_WEEKLY)
	a.WipeTournaments(ctx, api.TournamentInterval_MONTHLY)
}

func main() {
	app := NewWipeTournamentApp()
	lambda.Start(app.handler)
}

// WipeTournaments wipes all tournaments before the current start date
func (a *WipeTournamentApp) WipeTournaments(ctx context.Context, interval api.TournamentInterval) []model.RankedTournament {
	tournamentCurrentStartDate := a.tournamentService.GetTournamentStartDate(time.Now(), interval)
	var totalTournamentUsers []model.RankedTournament
	// Wipe tournaments before the current start date, loop until all tournaments are wiped
	limit := 100
	offset := 0
	maxLoops := 100
	countLoops := 0
	for countLoops < maxLoops {
		tournamentUsers, err := a.tournamentService.Database.GetTournamentsBeforeWipe(ctx, model.GetTournamentsBeforeWipeParams{
			TournamentStartedAt: tournamentCurrentStartDate,
			TournamentInterval:  model.TournamentTournamentInterval(interval.String()),
			Limit:               int32(limit),
			Offset:              int32(offset),
		})
		if err != nil {
			log.Fatalf("failed to get %s tournaments: %v", interval, err)
		}
		if len(tournamentUsers) == 0 {
			break
		}
		totalTournamentUsers = append(totalTournamentUsers, tournamentUsers...)
		deleteResult, err := a.tournamentService.Database.WipeTournaments(ctx, model.WipeTournamentsParams{
			TournamentStartedAt: tournamentCurrentStartDate,
			TournamentInterval:  model.TournamentTournamentInterval(interval.String()),
		})
		if err != nil {
			log.Fatalf("failed to wipe %s tournaments: %v", interval, err)
		}
		log.Printf("%s tournaments: %d, deleted: %d", interval, len(tournamentUsers), deleteResult)
		offset += limit
		countLoops++
	}
	if countLoops == maxLoops {
		log.Printf("max loops reached, %s tournaments: %d", interval, len(totalTournamentUsers))
	}
	return totalTournamentUsers
}

func stringToIntPanicOnError(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
