package sendEndedTournamentToThirdParty

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/sendEndedTournamentToThirdParty/model"
	"github.com/MorhafAlshibly/coanda/pkg/tournament"
)

type App struct {
	database                *model.Queries
	sql                     *sql.DB
	dailyTournamentMinute   uint16
	weeklyTournamentMinute  uint16
	weeklyTournamentDay     time.Weekday
	monthlyTournamentMinute uint16
	monthlyTournamentDay    uint8
	thirdPartyUri           string
	apiKeyHeader            string
	apiKey                  string
	topLimit                int32
	limit                   int32
}

func WithDatabase(database *model.Queries) func(*App) {
	return func(input *App) {
		input.database = database
	}
}

func WithSql(sql *sql.DB) func(*App) {
	return func(input *App) {
		input.sql = sql
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

func WithThirdPartyUri(thirdPartyUri string) func(*App) {
	return func(input *App) {
		input.thirdPartyUri = thirdPartyUri
	}
}

func WithApiKeyHeader(apiKeyHeader string) func(*App) {
	return func(input *App) {
		input.apiKeyHeader = apiKeyHeader
	}
}

func WithApiKey(apiKey string) func(*App) {
	return func(input *App) {
		input.apiKey = apiKey
	}
}

func WithTopLimit(topLimit int32) func(*App) {
	return func(input *App) {
		input.topLimit = topLimit
	}
}

func WithLimit(limit int32) func(*App) {
	return func(input *App) {
		input.limit = limit
	}
}

func NewApp(opts ...func(*App)) *App {
	app := App{
		apiKeyHeader: "x-api-key",
		topLimit:     1000,
		limit:        100,
	}
	for _, opt := range opts {
		opt(&app)
	}
	return &app
}

func (a *App) Handler(ctx context.Context) error {
	err := a.sendEndedTournamentsToThirdParty(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) sendEndedTournamentsToThirdParty(ctx context.Context) error {
	limit := a.limit
	offset := int32(0)
	wipeTimes := tournament.WipeTimes{
		DailyTournamentMinute:   a.dailyTournamentMinute,
		WeeklyTournamentMinute:  a.weeklyTournamentMinute,
		WeeklyTournamentDay:     a.weeklyTournamentDay,
		MonthlyTournamentMinute: a.monthlyTournamentMinute,
		MonthlyTournamentDay:    a.monthlyTournamentDay,
	}
	for {
		dailyTournaments, err := a.database.GetEndedTournaments(ctx, model.GetEndedTournamentsParams{
			TournamentStartedAt: tournament.GetStartTime(time.Now().UTC(), api.TournamentInterval_DAILY, wipeTimes),
			TournamentInterval:  model.TournamentTournamentIntervalDaily,
			Limit:               limit,
			Offset:              offset,
		})
		if err != nil {
			return err
		}
		weeklyTournaments, err := a.database.GetEndedTournaments(ctx, model.GetEndedTournamentsParams{
			TournamentStartedAt: tournament.GetStartTime(time.Now().UTC(), api.TournamentInterval_WEEKLY, wipeTimes),
			TournamentInterval:  model.TournamentTournamentIntervalWeekly,
			Limit:               limit,
			Offset:              offset,
		})
		if err != nil {
			return err
		}
		monthlyTournaments, err := a.database.GetEndedTournaments(ctx, model.GetEndedTournamentsParams{
			TournamentStartedAt: tournament.GetStartTime(time.Now().UTC(), api.TournamentInterval_MONTHLY, wipeTimes),
			TournamentInterval:  model.TournamentTournamentIntervalMonthly,
			Limit:               limit,
			Offset:              offset,
		})
		if err != nil {
			return err
		}
		tournaments := append(append(dailyTournaments, weeklyTournaments...), monthlyTournaments...)
		if len(tournaments) == 0 {
			break
		}
		for _, t := range tournaments {
			err = a.sendEndedTournamentUsersToThirdParty(ctx, t.Name, t.TournamentStartedAt, t.TournamentInterval)
			if err != nil {
				return err
			}
		}
		offset += limit
	}
	return nil
}

func (a *App) sendEndedTournamentUsersToThirdParty(ctx context.Context, name string, tournamentStartedAt time.Time, tournamentInterval model.TournamentTournamentInterval) error {
	tx, err := a.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := a.database.WithTx(tx)
	tournamentUsers, err := qtx.GetEndedTournamentUsers(ctx, model.GetEndedTournamentUsersParams{
		Name:                name,
		TournamentStartedAt: tournamentStartedAt,
		TournamentInterval:  tournamentInterval,
		Ranking:             uint64(a.topLimit),
	})
	if err != nil {
		return err
	}
	// Update tournament status to sent
	result, err := qtx.UpdateTournamentSentToThirdParty(ctx, model.UpdateTournamentSentToThirdPartyParams{
		Name:                name,
		TournamentStartedAt: tournamentStartedAt,
		TournamentInterval:  tournamentInterval,
	})
	if err != nil {
		return err
	}
	// Check if tournament was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return nil
	}
	// Create json and send to third party
	tournamentUsersData := []map[string]interface{}{}
	for _, tu := range tournamentUsers {
		tournamentUsersData = append(tournamentUsersData, map[string]interface{}{
			"id":                    tu.ID,
			"name":                  tu.Name,
			"tournament_interval":   tu.TournamentInterval,
			"user_id":               tu.UserID,
			"score":                 tu.Score,
			"ranking":               tu.Ranking,
			"data":                  string(tu.Data),
			"tournament_started_at": tu.TournamentStartedAt.Format(time.RFC3339),
			"created_at":            tu.CreatedAt.Format(time.RFC3339),
			"updated_at":            tu.UpdatedAt.Format(time.RFC3339),
		})
	}
	marshalledTournamentData, err := json.Marshal(tournamentUsersData)
	if err != nil {
		return err
	}
	// Send to third party with header
	thirdPartyUriWithID := fmt.Sprintf("%s?id=tournament-%s-%s-%s", a.thirdPartyUri, name, tournamentStartedAt.Format(time.RFC3339), tournamentInterval)
	req, err := http.NewRequest("POST", thirdPartyUriWithID, bytes.NewBuffer([]byte(fmt.Sprintf("{value: \"%s\"}", string(marshalledTournamentData)))))
	if err != nil {
		return err
	}
	req.Header.Set(a.apiKeyHeader, a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check response
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
