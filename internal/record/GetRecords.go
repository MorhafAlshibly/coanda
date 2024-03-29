package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
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
	max, page := pkg.ParsePagination(c.In.Max, c.In.Page, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	pipelineWithMatch := pipeline
	if c.In.Name != nil {
		if len(*c.In.Name) < int(c.service.minRecordNameLength) {
			c.Out = &api.GetRecordsResponse{
				Success: false,
				Records: nil,
				Error:   api.GetRecordsResponse_NAME_TOO_SHORT,
			}
			return nil
		}
		if len(*c.In.Name) > int(c.service.maxRecordNameLength) {
			c.Out = &api.GetRecordsResponse{
				Success: false,
				Records: nil,
				Error:   api.GetRecordsResponse_NAME_TOO_LONG,
			}
			return nil
		}
		pipelineWithMatch = mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "name", Value: c.In.Name},
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
	records, err := pkg.CursorToDocuments(ctx, cursor, toRecord, page, max)
	if err != nil {
		return err
	}
	c.Out = &api.GetRecordsResponse{
		Success: true,
		Records: records,
	}
	return nil
}
