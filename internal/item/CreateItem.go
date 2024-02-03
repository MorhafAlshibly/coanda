package item

import (
	"context"
	"database/sql"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"

	// // "github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/google/uuid"
)

type CreateItemCommand struct {
	service *Service
	In      *api.CreateItemRequest
	Out     *api.CreateItemResponse
}

func NewCreateItemCommand(service *Service, in *api.CreateItemRequest) *CreateItemCommand {
	return &CreateItemCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateItemCommand) Execute(ctx context.Context) error {
	// If the item has an expiry, add it to the map
	var expiresAt sql.NullTime
	if c.In.ExpiresAt != nil {
		// Check if the expiry is valid rfc3339
		expiresAtTime, err := time.Parse(time.RFC3339, *c.In.ExpiresAt)
		if err != nil {
			c.Out = &api.CreateItemResponse{
				Success: false,
				Error:   api.CreateItemResponse_EXPIRE_INVALID,
			}
			return nil
		}
		expiresAt = sql.NullTime{
			Time:  expiresAtTime,
			Valid: true,
		}
	} else {
		c.In.ExpiresAt = new(string)
		*c.In.ExpiresAt = ""
		expiresAt = sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	id := uuid.New().String()
	item := &api.Item{
		Id:        id,
		Data:      c.In.Data,
		ExpiresAt: *c.In.ExpiresAt,
	}
	// Add the item to the store
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	_, err = c.service.database.CreateItem(ctx, sqlc.CreateItemParams{
		ID:        id,
		Data:      data,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}
	// Allot the output
	c.Out = &api.CreateItemResponse{
		Success: true,
		Item:    item,
	}
	return nil
}
