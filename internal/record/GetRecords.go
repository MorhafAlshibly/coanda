package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type GetRecordsCommand struct {
	service *Service
	In      *api.GetRecordsRequest
	Out     *api.GetRecordsResponse
}

func NewGetRecordsCommand(service *Service, in *api.GetRecordsRequest) *GetRecordsCommand {
	return &GetRecordsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetRecordsCommand) Execute(ctx context.Context) error {
	if c.In.Max == nil {
		c.In.Max = new(uint64)
		*c.In.Max = c.service.defaultMaxPageLength
	}
	if c.In.Page == nil {
		c.In.Page = new(uint64)
		*c.In.Page = 1
	}
	pipelineWithMatch := pipeline
	if c.In.Name != nil {
		pipelineWithMatch = append(pipelineWithMatch, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "name", Value: *c.In.Name},
			}},
		})
	}
	cursor, err := c.service.db.Aggregate(ctx, pipelineWithMatch)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	records, err := toRecords(ctx, cursor, *c.In.Page, *c.In.Max)
	if err != nil {
		return err
	}
	c.Out = &api.GetRecordsResponse{
		Success: true,
		Records: records,
	}
	return nil
}
