package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"

	// "github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
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
	max := validation.ValidateMaxPageLength(c.In.Max, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	offset := conversion.PageToOffset(c.In.Page, max)
	var result []sqlc.RankedRecord
	var err error
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
		result, err = c.service.database.GetRecordsByName(ctx, sqlc.GetRecordsByNameParams{
			Name:   *c.In.Name,
			Offset: offset,
			Limit:  int32(max),
		})
	} else {
		result, err = c.service.database.GetRecords(ctx, sqlc.GetRecordsParams{
			Offset: offset,
			Limit:  int32(max),
		})
	}
	if err != nil {
		return err
	}
	records := make([]*api.Record, len(result))
	for i, record := range result {
		records[i], err = UnmarshalRecord(&record)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetRecordsResponse{
		Success: true,
		Records: records,
	}
	return nil
}
