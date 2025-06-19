package item

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetItemCommand struct {
	service *Service
	In      *api.ItemRequest
	Out     *api.GetItemResponse
}

func NewGetItemCommand(service *Service, in *api.ItemRequest) *GetItemCommand {
	return &GetItemCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForItemRequestError(c.In)
	if iErr != nil {
		c.Out = &api.GetItemResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.GetItemResponse_Error_value, api.GetItemResponse_ID_REQUIRED),
		}
		return nil
	}
	result, err := c.service.database.GetItem(ctx, model.GetItemParams{
		ID:   c.In.Id,
		Type: c.In.Type,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetItemResponse{
				Success: false,
				Item:    nil,
				Error:   api.GetItemResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	item, err := unmarshalItem(&result)
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
