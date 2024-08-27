package archiveEvent

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/archiveEvent/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type App struct {
	sql                 *sql.DB
	database            *model.Queries
	storage             storage.Storer
	folderPath          string
	sentToThirdPartyKey string
	limit               int32
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

func WithStorage(storage storage.Storer) func(*App) {
	return func(input *App) {
		input.storage = storage
	}
}

func WithFolderPath(folderPath string) func(*App) {
	return func(input *App) {
		input.folderPath = folderPath
	}
}

func WithSentToThirdPartyKey(sentToThirdPartyKey string) func(*App) {
	return func(input *App) {
		input.sentToThirdPartyKey = sentToThirdPartyKey
	}
}

func WithLimit(limit int32) func(*App) {
	return func(input *App) {
		input.limit = limit
	}
}

func NewApp(opts ...func(*App)) *App {
	app := App{
		folderPath: "/archive/events",
		limit:      100,
	}
	for _, opt := range opts {
		opt(&app)
	}
	return &app
}

func (a *App) Handler(ctx context.Context) error {
	// Folder path is the path where the events will be stored with current date and time in RFC3339 format
	folderPath := a.folderPath + "/" + time.Now().Format(time.RFC3339)
	err := a.archiveEvents(ctx, folderPath)
	if err != nil {
		fmt.Printf("failed to archive events: %v", err)
		return err
	}
	err = a.archiveEventUser(ctx, folderPath)
	if err != nil {
		fmt.Printf("failed to archive event users: %v", err)
		return err
	}
	err = a.archiveEventRound(ctx, folderPath)
	if err != nil {
		fmt.Printf("failed to archive event rounds: %v", err)
		return err
	}
	err = a.archiveEventRoundUser(ctx, folderPath)
	if err != nil {
		fmt.Printf("failed to archive event round users: %v", err)
		return err
	}
	return nil
}

func (a *App) archiveEvents(ctx context.Context, folderPath string) error {
	limit := a.limit
	offset := int32(0)
	for {
		tx, err := a.sql.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		qtx := a.database.WithTx(tx)
		events, err := qtx.GetEndedEvents(ctx, model.GetEndedEventsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		if len(events) == 0 {
			break
		}
		// Convert events to array of maps
		eventMaps := make([]map[string]interface{}, 0, len(events))
		for _, event := range events {
			eventMaps = append(eventMaps, map[string]interface{}{
				"id":         event.ID,
				"name":       event.Name,
				"data":       event.Data,
				"started_at": event.StartedAt,
				"created_at": event.CreatedAt,
				"updated_at": event.UpdatedAt,
			})
		}
		// Create a CSV file with the events
		csv, err := conversion.ArrayOfMapsToCsv(eventMaps)
		if err != nil {
			return err
		}
		// Gzip the CSV file
		var compressedCSV bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedCSV)
		if _, err := gzipWriter.Write([]byte(csv)); err != nil {
			return err
		}
		if err := gzipWriter.Close(); err != nil {
			return err
		}
		key := fmt.Sprintf("%s/event/%d-%d.csv.gz", folderPath, events[0].ID, events[len(events)-1].ID)
		// Store the compressed CSV file in the storage
		if err := a.storage.Store(ctx, key, compressedCSV.Bytes(), map[string]*string{
			a.sentToThirdPartyKey: nil,
		}); err != nil {
			return err
		}
		// Delete the events from the database
		_, err = qtx.DeleteEndedEvents(ctx, model.DeleteEndedEventsParams{
			MinID: events[0].ID,
			MaxID: events[len(events)-1].ID,
		})
		if err != nil {
			// If we fail to delete the events from the database, we should delete the stored CSV file
			_ = a.storage.Delete(ctx, key)
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		// Continue to the next batch of events
		offset += limit
	}
	return nil
}

func (a *App) archiveEventUser(ctx context.Context, folderPath string) error {
	limit := a.limit
	offset := int32(0)
	for {
		tx, err := a.sql.BeginTx(ctx, nil)
		if err != nil {
			fmt.Printf("failed to begin transaction: %v", err)
			return err
		}
		defer tx.Rollback()
		qtx := a.database.WithTx(tx)
		eventUsers, err := qtx.GetEndedEventUsers(ctx, model.GetEndedEventUsersParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			fmt.Printf("failed to get ended event users: %v", err)
			return err
		}
		if len(eventUsers) == 0 {
			break
		}
		// Convert event users to array of maps
		eventUserMaps := make([]map[string]interface{}, 0, len(eventUsers))
		for _, eventUser := range eventUsers {
			eventUserMaps = append(eventUserMaps, map[string]interface{}{
				"id":         eventUser.ID,
				"event_id":   eventUser.EventID,
				"user_id":    eventUser.UserID,
				"data":       eventUser.Data,
				"created_at": eventUser.CreatedAt,
				"updated_at": eventUser.UpdatedAt,
			})
		}
		// Create a CSV file with the event users
		csv, err := conversion.ArrayOfMapsToCsv(eventUserMaps)
		if err != nil {
			fmt.Printf("failed to convert event users to csv: %v", err)
			return err
		}
		// Gzip the CSV file
		var compressedCSV bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedCSV)
		if _, err := gzipWriter.Write([]byte(csv)); err != nil {
			fmt.Printf("failed to write compressed csv: %v", err)
			return err
		}
		if err := gzipWriter.Close(); err != nil {
			fmt.Printf("failed to close gzip writer: %v", err)
			return err
		}
		key := fmt.Sprintf("%s/event_user/%d-%d.csv.gz", folderPath, eventUsers[0].ID, eventUsers[len(eventUsers)-1].ID)
		// Store the compressed CSV file in the storage
		if err := a.storage.Store(ctx, key, compressedCSV.Bytes(), map[string]*string{
			a.sentToThirdPartyKey: nil,
		}); err != nil {
			fmt.Printf("failed to store compressed csv: %v", err)
			return err
		}
		// Delete the event users from the database
		_, err = qtx.DeleteEndedEventUsers(ctx, model.DeleteEndedEventUsersParams{
			MinID: eventUsers[0].ID,
			MaxID: eventUsers[len(eventUsers)-1].ID,
		})
		if err != nil {
			// If we fail to delete the event users from the database, we should delete the stored CSV file
			_ = a.storage.Delete(ctx, key)
			fmt.Printf("failed to delete ended event users: %v", err)
			return err
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf("failed to commit transaction: %v", err)
			return err
		}
		// Continue to the next batch of event users
		offset += limit
	}
	return nil
}

func (a *App) archiveEventRound(ctx context.Context, folderPath string) error {
	limit := a.limit
	offset := int32(0)
	for {
		tx, err := a.sql.BeginTx(ctx, nil)
		if err != nil {
			fmt.Printf("failed to begin transaction: %v", err)
			return err
		}
		defer tx.Rollback()
		qtx := a.database.WithTx(tx)
		eventRounds, err := qtx.GetEndedEventRounds(ctx, model.GetEndedEventRoundsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			fmt.Printf("failed to get ended event rounds: %v", err)
			return err
		}
		if len(eventRounds) == 0 {
			break
		}
		// Convert event rounds to array of maps
		eventRoundMaps := make([]map[string]interface{}, 0, len(eventRounds))
		for _, eventRound := range eventRounds {
			eventRoundMaps = append(eventRoundMaps, map[string]interface{}{
				"id":         eventRound.ID,
				"event_id":   eventRound.EventID,
				"name":       eventRound.Name,
				"scoring":    eventRound.Scoring,
				"data":       eventRound.Data,
				"ended_at":   eventRound.EndedAt,
				"created_at": eventRound.CreatedAt,
				"updated_at": eventRound.UpdatedAt,
			})
		}
		// Create a CSV file with the event rounds
		csv, err := conversion.ArrayOfMapsToCsv(eventRoundMaps)
		if err != nil {
			fmt.Printf("failed to convert event rounds to csv: %v", err)
			return err
		}
		// Gzip the CSV file
		var compressedCSV bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedCSV)
		if _, err := gzipWriter.Write([]byte(csv)); err != nil {
			fmt.Printf("failed to write compressed csv: %v", err)
			return err
		}
		if err := gzipWriter.Close(); err != nil {
			fmt.Printf("failed to close gzip writer: %v", err)
			return err
		}
		key := fmt.Sprintf("%s/event_round/%d-%d.csv.gz", folderPath, eventRounds[0].ID, eventRounds[len(eventRounds)-1].ID)
		// Store the compressed CSV file in the storage
		if err := a.storage.Store(ctx, key, compressedCSV.Bytes(), map[string]*string{
			a.sentToThirdPartyKey: nil,
		}); err != nil {
			fmt.Printf("failed to store compressed csv: %v", err)
			return err
		}
		// Delete the event rounds from the database
		_, err = qtx.DeleteEndedEventRounds(ctx, model.DeleteEndedEventRoundsParams{
			MinID: eventRounds[0].ID,
			MaxID: eventRounds[len(eventRounds)-1].ID,
		})
		if err != nil {
			// If we fail to delete the event rounds from the database, we should delete the stored CSV file
			_ = a.storage.Delete(ctx, key)
			fmt.Printf("failed to delete ended event rounds: %v", err)
			return err
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf("failed to commit transaction: %v", err)
			return err
		}
		// Continue to the next batch of event rounds
		offset += limit
	}
	return nil
}

func (a *App) archiveEventRoundUser(ctx context.Context, folderPath string) error {
	limit := a.limit
	offset := int32(0)
	for {
		tx, err := a.sql.BeginTx(ctx, nil)
		if err != nil {
			fmt.Printf("failed to begin transaction: %v", err)
			return err
		}
		defer tx.Rollback()
		qtx := a.database.WithTx(tx)
		eventRoundUsers, err := qtx.GetEndedEventRoundUsers(ctx, model.GetEndedEventRoundUsersParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			fmt.Printf("failed to get ended event round users: %v", err)
			return err
		}
		if len(eventRoundUsers) == 0 {
			break
		}
		// Convert event round users to array of maps
		eventRoundUserMaps := make([]map[string]interface{}, 0, len(eventRoundUsers))
		for _, eventRoundUser := range eventRoundUsers {
			eventRoundUserMaps = append(eventRoundUserMaps, map[string]interface{}{
				"id":             eventRoundUser.ID,
				"event_user_id":  eventRoundUser.EventUserID,
				"event_round_id": eventRoundUser.EventRoundID,
				"result":         eventRoundUser.Result,
				"data":           eventRoundUser.Data,
				"created_at":     eventRoundUser.CreatedAt,
				"updated_at":     eventRoundUser.UpdatedAt,
			})
		}
		// Create a CSV file with the event round users
		csv, err := conversion.ArrayOfMapsToCsv(eventRoundUserMaps)
		if err != nil {
			fmt.Printf("failed to convert event round users to csv: %v", err)
			return err
		}
		// Gzip the CSV file
		var compressedCSV bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedCSV)
		if _, err := gzipWriter.Write([]byte(csv)); err != nil {
			fmt.Printf("failed to write compressed csv: %v", err)
			return err
		}
		if err := gzipWriter.Close(); err != nil {
			fmt.Printf("failed to close gzip writer: %v", err)
			return err
		}
		key := fmt.Sprintf("%s/event_round_user/%d-%d.csv.gz", folderPath, eventRoundUsers[0].ID, eventRoundUsers[len(eventRoundUsers)-1].ID)
		// Store the compressed CSV file in the storage
		if err := a.storage.Store(ctx, key, compressedCSV.Bytes(), map[string]*string{
			a.sentToThirdPartyKey: nil,
		}); err != nil {
			fmt.Printf("failed to store compressed csv: %v", err)
			return err
		}
		// Delete the event round users from the database
		_, err = qtx.DeleteEndedEventRoundUsers(ctx, model.DeleteEndedEventRoundUsersParams{
			MinID: eventRoundUsers[0].ID,
			MaxID: eventRoundUsers[len(eventRoundUsers)-1].ID,
		})
		if err != nil {
			// If we fail to delete the event round users from the database, we should delete the stored CSV file
			_ = a.storage.Delete(ctx, key)
			fmt.Printf("failed to delete ended event round users: %v", err)
			return err
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf("failed to commit transaction: %v", err)
			return err
		}
		// Continue to the next batch of event round users
		offset += limit
	}
	return nil
}
