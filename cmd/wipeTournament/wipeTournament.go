package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	lambdaFunc "github.com/aws/aws-lambda-go/lambda"
	"github.com/peterbourgon/ff/v4"
	"github.com/robfig/cron/v3"
)

var (
	fs                      = ff.NewFlagSet("wipeTournament")
	lambda                  = fs.BoolLong("lambda", "if running as a lambda function")
	cronSchedule            = fs.StringLong("cronSchedule", "* * * * *", "the cron schedule to run the handler")
	dsn                     = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	dailyTournamentMinute   = fs.UintLong("dailyTournamentMinute", 0, "the minute of the day to start the daily tournament")
	weeklyTournamentMinute  = fs.UintLong("weeklyTournamentMinute", 0, "the minute of the week to start the weekly tournament")
	weeklyTournamentDay     = fs.UintLong("weeklyTournamentDay", 0, "the day of the week to start the weekly tournament")
	monthlyTournamentMinute = fs.UintLong("monthlyTournamentMinute", 0, "the minute of the month to start the monthly tournament")
	monthlyTournamentDay    = fs.UintLong("monthlyTournamentDay", 1, "the day of the month to start the monthly tournament")
)

func main() {
	app := NewWipeTournamentApp()
	if !*lambda {
		// Run the handler on a cron job
		c := cron.New()
		c.AddFunc(*cronSchedule, func() {
			if err := app.handler(context.Background()); err != nil {
				log.Fatalf("failed to run handler: %v", err)
			}
		})
		c.Start()
		return
	}
	// Run the lambda if not running on a cron job
	lambdaFunc.Start(app.handler)
}

type WipeTournamentApp struct {
	tournamentService *tournament.Service
}

func NewWipeTournamentApp() *WipeTournamentApp {
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("WIPE_TOURNAMENT"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}
	dbConn, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db := model.New(dbConn)
	return &WipeTournamentApp{
		tournamentService: tournament.NewService(
			tournament.WithSql(dbConn),
			tournament.WithDatabase(db),
			tournament.WithDailyTournamentMinute(uint16(*dailyTournamentMinute)),
			tournament.WithWeeklyTournamentMinute(uint16(*weeklyTournamentMinute)),
			tournament.WithWeeklyTournamentDay(time.Weekday(*weeklyTournamentDay)),
			tournament.WithMonthlyTournamentMinute(uint16(*monthlyTournamentMinute)),
			tournament.WithMonthlyTournamentDay(uint8(*monthlyTournamentDay)),
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
