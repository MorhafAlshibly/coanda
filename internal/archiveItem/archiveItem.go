package archiveItem

import (
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/archiveItem/model"
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

func WithLimit(limit int) func(*App) {
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
	// Start transaction
	// Create a new folder in the archive for the items named with the current date
	// Get all items that have expired (sorted by id ascending) in a loop with a set limit so we don't run out of memory
	// Create a CSV file with the items, gzip it and store it in the folder named by minimum id and maximum id
	// Continue the loop until there are no more items
	// Delete the items from the database
	// Commit transaction
	tx, err := a.sql.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()
	qtx := a.database.WithTx(tx)
	// Folder path for the current date and time in RFC3339 format
	folderPath := a.folderPath + time.Now().Format(time.RFC3339)

}

func (a *App) archiveItems(ctx context.Context, folderPath string) error {
	// Get all items that have expired (sorted by id ascending) in a loop with a set limit so we don't run out of memory
	limit := a.limit
	offset := int32(0)
	for {
		items, err := a.database.GetExpiredItems(ctx, model.GetExpiredItemsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		if len(items) == 0 {
			break
		}
		// Create a CSV file with the items, gzip it and store it in the folder named by minimum id and maximum id
		csv, err := createCompressedCSV(items)
		if err != nil {
			return err
		}
		err = a.storage.Store(ctx, folderPath+"/"+items[0].ID+"-"+items[len(items)-1].ID+".csv.gz", csv)
		// Continue the loop until there are no more items
	}
}

func createCompressedCSV(items []model.Item) ([]byte, error) {
	// Create a CSV file with the items and gzip it
	csvString := "id,type,data,created_at,updated_at,expires_at"
	for _, item := range items {
		csvString += "\n" + item.ID + "," + item.Type + "," + string(item.Data) + "," + item.CreatedAt.Format(time.RFC3339) + "," + item.UpdatedAt.Format(time.RFC3339)
		if item.ExpiresAt.Valid {
			csvString += "," + item.ExpiresAt.Time.Format(time.RFC3339)
		}
	}
	// Gzip the CSV file
	var compressedCSV []byte
	gz := gzip.NewWriter(&compressedCSV)
	_, err := gz.Write([]byte(csvString))
	if err != nil {
		return err
	}
	// Return the compressed CSV file
	return nil
}
