package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
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
	if len(c.In.Query) < int(c.service.minTeamNameLength) {
		c.Out = &api.SearchTeamsResponse{
			Success: false,
			Teams:   nil,
			Error:   api.SearchTeamsResponse_QUERY_TOO_SHORT,
		}
		return nil
	}
	max, page := pkg.ParsePagination(c.In.Pagination.Max, c.In.Pagination.Page, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	cursor, err := c.service.db.Aggregate(ctx, append(pipeline, searchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	teams, err := pkg.CursorToDocuments(ctx, cursor, toTeam, page, max)
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
