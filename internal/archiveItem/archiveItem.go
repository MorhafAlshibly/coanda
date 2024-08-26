package archiveItem

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/archiveItem/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type App struct {
	sql        *sql.DB
	database   *model.Queries
	storage    storage.Storer
	folderPath string
	limit      int32
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

func WithLimit(limit int32) func(*App) {
	return func(input *App) {
		input.limit = limit
	}
}

func NewApp(opts ...func(*App)) *App {
	app := App{
		folderPath: "/archive/items",
		limit:      100,
	}
	for _, opt := range opts {
		opt(&app)
	}
	return &app
}

func (a *App) Handler(ctx context.Context) error {
	// Folder path for the current date and time in RFC3339 format
	folderPath := a.folderPath + time.Now().Format(time.RFC3339)
	err := a.archiveItems(ctx, folderPath)
	if err != nil {
		fmt.Printf("failed to archive items: %v", err)
		return err
	}
	return nil
}

func (a *App) archiveItems(ctx context.Context, folderPath string) error {
	// Get all items that have expired (sorted by id ascending) in a loop with a set limit so we don't run out of memory
	limit := a.limit
	offset := int32(0)
	for {
		tx, err := a.sql.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		qtx := a.database.WithTx(tx)
		items, err := qtx.GetExpiredItems(ctx, model.GetExpiredItemsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		if len(items) == 0 {
			break
		}
		// Convert items to array of maps
		itemsMap := make([]map[string]interface{}, len(items))
		for i, item := range items {
			itemsMap[i] = map[string]interface{}{
				"id":         item.ID,
				"type":       item.Type,
				"data":       item.Data,
				"created_at": item.CreatedAt,
				"updated_at": item.UpdatedAt,
				"expires_at": item.ExpiresAt,
			}
		}
		// Create a CSV file with the items, gzip it and store it in the folder named by minimum id and maximum id
		csv, err := conversion.ArrayOfMapsToCSV(itemsMap, nil)
		if err != nil {
			return err
		}
		// Gzip the CSV file
		var compressedCSV bytes.Buffer
		gz := gzip.NewWriter(&compressedCSV)
		_, err = gz.Write([]byte(csv))
		if err != nil {
			return err
		}
		err = gz.Close()
		if err != nil {
			return err
		}
		// Store the compressed CSV file in the storage
		err = a.storage.Store(ctx, folderPath+"/"+items[0].ID+"-"+items[len(items)-1].ID+".csv.gz", compressedCSV.Bytes())
		if err != nil {
			return err
		}
		// Delete the items from the database
		_, err = qtx.DeleteExpiredItems(ctx, model.DeleteExpiredItemsParams{
			MinID: items[0].ID,
			MaxID: items[len(items)-1].ID,
		})
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		// Continue the loop until there are no more items
		offset += int32(len(items))
	}
	return nil
}
