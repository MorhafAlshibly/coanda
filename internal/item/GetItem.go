package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetItemCommand struct {
	service *Service
	In      *api.GetItemRequest
	Out     *api.Item
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
	c.Out, err = objectToItem(object)
	if err != nil {
		return err
	}
	return nil
}
