package item

import (
	"context"

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
	object, err := c.service.store.Get(ctx, c.In.Id, c.In.Type)
	if err != nil {
		return err
	}
	item, err := objectToItem(object)
	if err != nil {
		c.Out = &api.GetItemResponse{
			Success: false,
			Item:    nil,
			Error:   api.GetItemResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.GetItemResponse{
		Success: true,
		Item:    item,
		Error:   api.GetItemResponse_NONE,
	}
	return nil
}
