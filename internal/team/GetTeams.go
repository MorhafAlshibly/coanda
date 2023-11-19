package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
)

type GetTeamsCommand struct {
	service *Service
	In      *api.GetTeamsRequest
	Out     *api.GetTeamsResponse
}

func NewGetTeamsCommand(service *Service, in *api.GetTeamsRequest) *GetTeamsCommand {
	return &GetTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	max, page := pkg.ParsePagination(c.In.Max, c.In.Page, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	cursor, err := c.service.db.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	teams, err := pkg.CursorToDocuments(ctx, cursor, toTeam, page, max)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamsResponse{
		Success: true,
		Teams:   teams,
	}
	return nil
}
