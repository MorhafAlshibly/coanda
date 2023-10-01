package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type SearchTeamsCommand struct {
	service *Service
	In      *api.SearchTeamsRequest
	Out     *api.SearchTeamsResponse
}

func NewSearchTeamsCommand(service *Service, in *api.SearchTeamsRequest) *SearchTeamsCommand {
	return &SearchTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *SearchTeamsCommand) Execute(ctx context.Context) error {
	searchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "$regex", Value: c.In.Query},
				{Key: "$options", Value: "i"},
			}},
		}},
	}
	if len(c.In.Query) < c.service.minTeamNameLength {
		c.Out = &api.SearchTeamsResponse{
			Success: false,
			Teams:   nil,
			Error:   api.SearchTeamsResponse_QUERY_TOO_SHORT,
		}
		return nil
	}
	if c.In.Pagination.Max == nil {
		c.In.Pagination.Max = new(uint64)
		*c.In.Pagination.Max = c.service.defaultMaxPageLength
	}
	if c.In.Pagination.Page == nil {
		c.In.Pagination.Page = new(uint64)
		*c.In.Pagination.Page = 1
	}
	cursor, err := c.service.db.Aggregate(ctx, append(pipeline, searchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	teams, err := toTeams(ctx, cursor, *c.In.Pagination.Page, *c.In.Pagination.Max)
	if err != nil {
		return err
	}
	c.Out = &api.SearchTeamsResponse{
		Success: true,
		Teams:   teams,
		Error:   api.SearchTeamsResponse_NONE,
	}
	return nil
}
