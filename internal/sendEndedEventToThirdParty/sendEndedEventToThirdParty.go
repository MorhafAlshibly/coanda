package sendEndedEventToThirdParty

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/sendEndedEventToThirdParty/model"
)

type App struct {
	database      *model.Queries
	sql           *sql.DB
	thirdPartyUri string
	apiKeyHeader  string
	apiKey        string
	topLimit      int32
	limit         int32
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
	err := a.sendEndedEventRoundsToThirdParty(ctx)
	if err != nil {
		fmt.Printf("failed to send ended event rounds to third party: %v", err)
		return err
	}
	err = a.sendEndedEventsToThirdParty(ctx)
	if err != nil {
		fmt.Printf("failed to send ended events to third party: %v", err)
		return err
	}
	return nil
}

func (a *App) sendEndedEventRoundsToThirdParty(ctx context.Context) error {
	limit := a.limit
	offset := int32(0)
	for {
		eventRounds, err := a.database.GetEndedEventRounds(ctx, model.GetEndedEventRoundsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		if len(eventRounds) == 0 {
			break
		}
		for _, eventRound := range eventRounds {
			err := a.sendEndedEventRoundToThirdParty(ctx, eventRound)
			if err != nil {
				return err
			}
		}
		offset += limit
	}
	return nil
}

func (a *App) sendEndedEventRoundToThirdParty(ctx context.Context, eventRound model.GetEndedEventRoundsRow) error {
	tx, err := a.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := a.database.WithTx(tx)
	leaderboard, err := qtx.GetEndedEventRoundLeaderboard(ctx, model.GetEndedEventRoundLeaderboardParams{
		EventRoundID: eventRound.ID,
		Ranking:      uint64(a.topLimit),
	})
	if err != nil {
		return err
	}
	// Update round status to sent
	result, err := qtx.UpdateEventRoundSentToThirdParty(ctx, eventRound.ID)
	if err != nil {
		return err
	}
	// Check if round was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return nil
	}
	// Create json and send to third party
	leaderboardData := []map[string]interface{}{}
	for _, lb := range leaderboard {
		leaderboardData = append(leaderboardData, map[string]interface{}{
			"id":             lb.ID,
			"event_user_id":  lb.EventUserID,
			"client_user_id": lb.ClientUserID,
			"result":         lb.Result,
			"score":          lb.Score,
			"ranking":        lb.Ranking,
			"data":           string(lb.Data),
			"created_at":     lb.CreatedAt.Format(time.RFC3339),
			"updated_at":     lb.UpdatedAt.Format(time.RFC3339),
		})
	}
	roundData := map[string]interface{}{
		"id":         eventRound.ID,
		"event_id":   eventRound.EventID,
		"scoring":    eventRound.Scoring,
		"users":      leaderboardData,
		"data":       string(eventRound.Data),
		"ended_at":   eventRound.EndedAt.Format(time.RFC3339),
		"created_at": eventRound.CreatedAt.Format(time.RFC3339),
		"updated_at": eventRound.UpdatedAt.Format(time.RFC3339),
	}
	marshalledRoundData, err := json.Marshal(roundData)
	if err != nil {
		return err
	}
	// Send to third party with header
	thirdPartyUriWithID := fmt.Sprintf("%s?id=event-%d-round-%d", a.thirdPartyUri, eventRound.EventID, eventRound.ID)
	req, err := http.NewRequest("POST", thirdPartyUriWithID, bytes.NewBuffer([]byte(fmt.Sprintf("{value: \"%s\"}", string(marshalledRoundData)))))
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

func (a *App) sendEndedEventsToThirdParty(ctx context.Context) error {
	limit := a.limit
	offset := int32(0)
	for {
		events, err := a.database.GetEndedEvents(ctx, model.GetEndedEventsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		if len(events) == 0 {
			break
		}
		for _, event := range events {
			err := a.sendEndedEventToThirdParty(ctx, event)
			if err != nil {
				return err
			}
		}
		offset += limit
	}
	return nil
}

func (a *App) sendEndedEventToThirdParty(ctx context.Context, event model.GetEndedEventsRow) error {
	tx, err := a.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := a.database.WithTx(tx)
	leaderboard, err := qtx.GetEndedEventLeaderboard(ctx, model.GetEndedEventLeaderboardParams{
		ID:      event.ID,
		Ranking: uint64(a.topLimit),
	})
	if err != nil {
		return err
	}
	// Update event status to sent
	result, err := qtx.UpdateEventSentToThirdParty(ctx, event.ID)
	if err != nil {
		return err
	}
	// Check if event was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return nil
	}
	// Create json and send to third party
	leaderboardData := []map[string]interface{}{}
	for _, lb := range leaderboard {
		leaderboardData = append(leaderboardData, map[string]interface{}{
			"id":             lb.ID,
			"event_id":       lb.EventID,
			"client_user_id": lb.ClientUserID,
			"score":          lb.Score,
			"ranking":        lb.Ranking,
			"data":           string(lb.Data),
			"created_at":     lb.CreatedAt.Format(time.RFC3339),
			"updated_at":     lb.UpdatedAt.Format(time.RFC3339),
		})
	}
	eventData := map[string]interface{}{
		"id":         event.ID,
		"name":       event.Name,
		"users":      leaderboardData,
		"data":       string(event.Data),
		"started_at": event.StartedAt.Format(time.RFC3339),
		"created_at": event.CreatedAt.Format(time.RFC3339),
		"updated_at": event.UpdatedAt.Format(time.RFC3339),
	}
	marshalledEventData, err := json.Marshal(eventData)
	if err != nil {
		return err
	}
	// Send to third party with header
	thirdPartyUriWithID := fmt.Sprintf("%s?id=event-%d", a.thirdPartyUri, event.ID)
	req, err := http.NewRequest("POST", thirdPartyUriWithID, bytes.NewBuffer([]byte(fmt.Sprintf("{value: \"%s\"}", string(marshalledEventData)))))
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
