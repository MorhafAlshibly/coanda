package item

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/bytedance/sonic"
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
	if c.In.Data == nil {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_DATA_NOT_SET,
		}
		return nil
	}
	marshalled, err := sonic.Marshal(c.In.Data)
	if err != nil {
		return err
	}
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
	mapData := map[string]string{
		"Type":   c.In.Type,
		"Data":   string(marshalled),
		"Expire": *c.In.Expire,
	}
	// Add the item to the store
	object, err := c.service.store.Add(ctx, c.In.Type, mapData)
	if err != nil {
		return err
	}
	// Allot the output
	c.Out = &api.CreateItemResponse{
		Success: true,
		Item: &api.Item{
			Id:     object.Key,
			Type:   c.In.Type,
			Data:   c.In.Data,
			Expire: *c.In.Expire,
		},
	}
	return nil
}
