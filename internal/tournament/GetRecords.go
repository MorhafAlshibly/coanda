package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
		if len(*c.In.Name) < c.service.minRecordNameLength {
			c.Out = &api.GetRecordsResponse{
				Success: false,
				Records: nil,
				Error:   api.GetRecordsResponse_NAME_TOO_SHORT,
			}
			return nil
		}
		pipelineWithMatch = mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "name", Value: *c.In.Name},
				}},
			},
		}
		for _, stage := range pipeline {
			pipelineWithMatch = append(pipelineWithMatch, stage)
		}
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
