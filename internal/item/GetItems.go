package item

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
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
	var outs []*api.Item
	// Check if the type is not nil, if it is, check if it is within the min and max length
	if c.In.Type != nil {
		if len(*c.In.Type) < int(c.service.minTypeLength) {
			c.Out = &api.GetItemsResponse{
				Success: false,
				Items:   nil,
				Error:   api.GetItemsResponse_TYPE_TOO_SHORT,
			}
			return nil
		}
		if len(*c.In.Type) > int(c.service.maxTypeLength) {
			c.Out = &api.GetItemsResponse{
				Success: false,
				Items:   nil,
				Error:   api.GetItemsResponse_TYPE_TOO_LONG,
			}
			return nil
		}
	}
	max := validation.ValidateMaxPageLength(c.In.Max, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	var items []map[string]any
	var err error
	// If the last evaluated id is not nil, set the exclusive start key to the last evaluated id
	var exclusiveStartKey map[string]string
	if c.In.Type == nil {
		if c.In.LastEvaluatedId != nil {
			exclusiveStartKey = map[string]string{
				"id": *c.In.LastEvaluatedId,
			}
		}
		items, err = c.service.database.Scan(ctx, &dynamoTable.ScanInput{
			ExclusiveStartKey: exclusiveStartKey,
			Max:               max,
		})
	} else {
		if c.In.LastEvaluatedId != nil {
			exclusiveStartKey = map[string]string{
				"id":   *c.In.LastEvaluatedId,
				"type": *c.In.Type,
			}
		}
		items, err = c.service.database.Query(ctx, &dynamoTable.QueryInput{
			ExclusiveStartKey:      exclusiveStartKey,
			KeyConditionExpression: "#type = :type",
			Max:                    max,
			ExpressionAttributeNames: map[string]string{
				"#type": "type",
			},
			ExpressionAttributeValues: map[string]any{
				":type": *c.In.Type,
			},
		})
	}
	if err != nil {
		return err
	}
	for _, object := range items {
		item, err := UnmarshalItem(object)
		if err != nil {
			return err
		}
		outs = append(outs, item)
	}
	c.Out = &api.GetItemsResponse{
		Success: true,
		Items:   outs,
		Error:   api.GetItemsResponse_NONE,
	}
	return nil
}
