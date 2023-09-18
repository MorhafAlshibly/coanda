package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type GetItemsCommand struct {
	service *Service
	In      *api.GetItemsRequest
	Out     []*api.Item
}

func NewGetItemsCommand(service *Service, in *api.GetItemsRequest) *GetItemsCommand {
	return &GetItemsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemsCommand) Execute(ctx context.Context) error {
	var items []*storage.Object
	var outs []*api.Item
	// If the type is not nil, set the filter to the type
	filter := ""
	if c.In.Type != "" {
		filter = "PartitionKey eq '" + c.In.Type + "'"
	}
	items, err := c.service.store.Query(ctx, filter, int32(c.In.Max), int(c.In.Page))
	if err != nil {
		return err
	}
	for _, item := range items {
		out, err := objectToItem(item)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	}
	return nil
}
