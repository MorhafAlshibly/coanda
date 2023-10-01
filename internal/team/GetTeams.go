package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
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
	if c.In.Max == nil {
		c.In.Max = new(uint64)
		*c.In.Max = c.service.defaultMaxPageLength
	}
	if c.In.Page == nil {
		c.In.Page = new(uint64)
		*c.In.Page = 1
	}
	cursor, err := c.service.db.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	teams, err := toTeams(ctx, cursor, *c.In.Page, *c.In.Max)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamsResponse{
		Success: true,
		Teams:   teams,
	}
	return nil
}
