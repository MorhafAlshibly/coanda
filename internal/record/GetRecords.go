package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
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
	}
	result, err := c.service.database.GetRecords(ctx, model.GetRecordsParams{
		Name:   validation.ValidateAnSqlNullString(c.In.Name),
		UserID: validation.ValidateAUint64ToSqlNullInt64(c.In.UserId),
		Limit:  limit,
		Offset: offset,
	})
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
