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
	max := uint8(c.In.Max)
	if max == 0 {
		max = c.service.defaultMaxPageLength
	}
	if max > c.service.maxMaxPageLength {
		max = c.service.maxMaxPageLength
	}
	if c.In.Page == 0 {
		c.In.Page = 1
	}
	cursor, err := c.service.db.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	teams, err := toTeams(ctx, cursor, c.In.Page, max)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamsResponse{
		Success: true,
		Teams:   teams,
	}
	return nil
}
