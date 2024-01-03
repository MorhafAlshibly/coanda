package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
)

type GetItemsCommand struct {
	service *Service
	In      *api.GetItemsRequest
	Out     *api.GetItemsResponse
}

func NewGetItemsCommand(service *Service, in *api.GetItemsRequest) *GetItemsCommand {
	return &GetItemsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetItemsCommand) Execute(ctx context.Context) error {
	max := validation.ValidateMaxPageLength(c.In.Max, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	offset := conversion.PageToOffset(c.In.Page, max)
	result, err := c.service.database.GetItems(ctx, sqlc.GetItemsParams{
		Offset: offset,
		Limit:  int32(max),
	})
	if err != nil {
		return err
	}
	items := make([]*api.Item, len(result))
	for i, item := range result {
		items[i], err = UnmarshalItem(&item)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetItemsResponse{
		Success: true,
		Items:   items,
	}
	return nil
}
