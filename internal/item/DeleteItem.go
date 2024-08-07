package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteItemCommand struct {
	service *Service
	In      *api.ItemRequest
	Out     *api.ItemResponse
}

func NewDeleteItemCommand(service *Service, in *api.ItemRequest) *DeleteItemCommand {
	return &DeleteItemCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteItemCommand) Execute(ctx context.Context) error {
	iErr := c.service.checkForItemRequestError(c.In)
	if iErr != nil {
		c.Out = &api.ItemResponse{
			Success: false,
			Error:   conversion.Enum(*iErr, api.ItemResponse_Error_value, api.ItemResponse_NOT_FOUND),
		}
		return nil
	}
	result, err := c.service.database.DeleteItem(ctx, model.DeleteItemParams{
		ID:   c.In.Id,
		Type: c.In.Type,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.ItemResponse{
			Success: false,
			Error:   api.ItemResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.ItemResponse{
		Success: true,
		Error:   api.ItemResponse_NONE,
	}
	return nil
}
