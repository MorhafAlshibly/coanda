package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
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
	if len(c.In.Type) < int(c.service.minTypeLength) {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_TYPE_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Type) > int(c.service.maxTypeLength) {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_TYPE_TOO_LONG,
		}
		return nil
	}
	// If the item has an expiry, add it to the map
	if c.In.Expire != nil {
		// Check if the expiry is valid rfc3339
		_, err := time.Parse(time.RFC3339, *c.In.Expire)
		if err != nil {
			c.Out = &api.CreateItemResponse{
				Success: false,
				Error:   api.CreateItemResponse_EXPIRE_INVALID,
			}
			return nil
		}
	} else {
		// If the item has no expiry, set it to empty string
		c.In.Expire = new(string)
		*c.In.Expire = ""
	}
	id := uuid.New().String()
	item := &api.Item{
		Id:     id,
		Type:   c.In.Type,
		Data:   c.In.Data,
		Expire: *c.In.Expire,
	}
	// Add the item to the store
	object, err := MarshalItem(item)
	if err != nil {
		return err
	}
	err = c.service.database.PutItem(ctx, &dynamoTable.PutItemInput{
		Item: object,
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
