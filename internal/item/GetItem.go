package item

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetItemCommand struct {
	service *Service
	In      *api.GetItemRequest
	Out     *api.GetItemResponse
}

func NewGetItemCommand(service *Service, in *api.GetItemRequest) *GetItemCommand {
	return &GetItemCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemCommand) Execute(ctx context.Context) error {
	if c.In.Id == "" {
		c.Out = &api.GetItemResponse{
			Success: false,
			Item:    nil,
			Error:   api.GetItemResponse_ID_NOT_SET,
		}
		return nil
	}
	result, err := c.service.database.GetItem(ctx, c.In.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetItemResponse{
				Success: false,
				Item:    nil,
				Error:   api.GetItemResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	item, err := UnmarshalItem(&result)
	if err != nil {
		return err
	}
	c.Out = &api.GetItemResponse{
		Success: true,
		Item:    item,
		Error:   api.GetItemResponse_NONE,
	}
	return nil
}
