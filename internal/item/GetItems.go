package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	result, err := c.service.database.GetItems(ctx, model.GetItemsParams{
		Type:   validation.ValidateAnSqlNullString(c.In.Type),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}
	items := make([]*api.Item, len(result))
	for i, item := range result {
		items[i], err = unmarshalItem(&item)
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
