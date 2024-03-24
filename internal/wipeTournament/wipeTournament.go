package wipeTournament

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
)

type App struct {
	tournamentService       *tournament.Service
	sql                     *sql.DB
	database                *model.Queries
	dailyTournamentMinute   uint16
	weeklyTournamentMinute  uint16
	weeklyTournamentDay     time.Weekday
	monthlyTournamentMinute uint16
	monthlyTournamentDay    uint8
}

func WithSql(sql *sql.DB) func(*App) {
	return func(input *App) {
		input.sql = sql
	}
}

func WithDatabase(database *model.Queries) func(*App) {
	return func(input *App) {
		input.database = database
	}
}

func WithDailyTournamentMinute(dailyTournamentMinute uint16) func(*App) {
	return func(input *App) {
		input.dailyTournamentMinute = dailyTournamentMinute
	}
}

func WithWeeklyTournamentMinute(weeklyTournamentMinute uint16) func(*App) {
	return func(input *App) {
		input.weeklyTournamentMinute = weeklyTournamentMinute
	}
}

func WithWeeklyTournamentDay(weeklyTournamentDay time.Weekday) func(*App) {
	return func(input *App) {
		input.weeklyTournamentDay = weeklyTournamentDay
	}
}

func WithMonthlyTournamentMinute(monthlyTournamentMinute uint16) func(*App) {
	return func(input *App) {
		input.monthlyTournamentMinute = monthlyTournamentMinute
	}
}

func WithMonthlyTournamentDay(monthlyTournamentDay uint8) func(*App) {
	return func(input *App) {
		input.monthlyTournamentDay = monthlyTournamentDay
	}
}

func NewApp(opts ...func(*App)) *App {
	app := App{
		dailyTournamentMinute:   0,
		weeklyTournamentMinute:  0,
		weeklyTournamentDay:     time.Monday,
		monthlyTournamentMinute: 0,
		monthlyTournamentDay:    1,
	}
	for _, opt := range opts {
		opt(&app)
	}
	app.tournamentService = tournament.NewService(
		tournament.WithSql(app.sql),
		tournament.WithDatabase(app.database),
		tournament.WithDailyTournamentMinute(app.dailyTournamentMinute),
		tournament.WithWeeklyTournamentMinute(app.weeklyTournamentMinute),
		tournament.WithWeeklyTournamentDay(app.weeklyTournamentDay),
		tournament.WithMonthlyTournamentMinute(app.monthlyTournamentMinute),
		tournament.WithMonthlyTournamentDay(app.monthlyTournamentDay),
	)
	return &app
}

func (a *App) Handler(ctx context.Context) error {
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
	dailyTournament, err := a.wipeTournaments(ctx, api.TournamentInterval_DAILY)
	if err != nil {
		log.Printf("failed to wipe daily tournaments: %v", err)
		return err
	}
	weeklyTournament, err := a.wipeTournaments(ctx, api.TournamentInterval_WEEKLY)
	if err != nil {
		log.Printf("failed to wipe weekly tournaments: %v", err)
		return err
	}
	monthlyTournament, err := a.wipeTournaments(ctx, api.TournamentInterval_MONTHLY)
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

// wipeTournaments wipes all tournaments before the current start date
func (a *App) wipeTournaments(ctx context.Context, interval api.TournamentInterval) (int64, error) {
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
