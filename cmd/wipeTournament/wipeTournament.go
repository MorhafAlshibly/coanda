package main

import (
	"context"
	"database/sql"
	"errors"
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

func (a *WipeTournamentApp) handler(ctx context.Context) error {
	sql := a.tournamentService.Sql()
	db := a.tournamentService.Database()
	tx, err := sql.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()
	qtx := db.WithTx(tx)
	a.tournamentService.SetDatabase(qtx)
	// Wipe all tournaments
	dailyTournament, err := a.WipeTournaments(ctx, api.TournamentInterval_DAILY)
	if err != nil {
		log.Printf("failed to wipe daily tournaments: %v", err)
		return err
	}
	weeklyTournament, err := a.WipeTournaments(ctx, api.TournamentInterval_WEEKLY)
	if err != nil {
		log.Printf("failed to wipe weekly tournaments: %v", err)
		return err
	}
	monthlyTournament, err := a.WipeTournaments(ctx, api.TournamentInterval_MONTHLY)
	if err != nil {
		log.Printf("failed to wipe monthly tournaments: %v", err)
		return err
	}
	log.Printf("wiped %d daily, %d weekly, and %d monthly tournaments", dailyTournament, weeklyTournament, monthlyTournament)
	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return nil
}

func main() {
	app := NewWipeTournamentApp()
	lambda.Start(app.handler)
}

// WipeTournaments wipes all tournaments before the current start date
func (a *WipeTournamentApp) WipeTournaments(ctx context.Context, interval api.TournamentInterval) (int64, error) {
	tournamentCurrentStartDate := a.tournamentService.GetTournamentStartDate(time.Now(), interval)
	// Wipe tournaments before the current start date
	result, err := a.tournamentService.Database().ArchiveTournaments(ctx, model.ArchiveTournamentsParams{
		TournamentStartedAt: tournamentCurrentStartDate,
		TournamentInterval:  model.TournamentTournamentInterval(interval.String()),
	})
	if err != nil {
		log.Printf("failed to archive %s tournaments: %v", interval, err)
		return 0, err
	}
	archiveRowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected: %v", err)
		return 0, err
	}
	result, err = a.tournamentService.Database().WipeTournaments(ctx, model.WipeTournamentsParams{
		TournamentStartedAt: tournamentCurrentStartDate,
		TournamentInterval:  model.TournamentTournamentInterval(interval.String()),
	})
	if err != nil {
		log.Printf("failed to delete %s tournaments: %v", interval, err)
		return 0, err
	}
	wipeRowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected: %v", err)
		return 0, err
	}
	if archiveRowsAffected != wipeRowsAffected {
		return 0, errors.New("archive rows affected not equal to wipe rows affected")
	}
	return archiveRowsAffected, nil
}

func stringToIntPanicOnError(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
