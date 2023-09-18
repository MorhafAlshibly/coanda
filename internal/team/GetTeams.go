package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetTeamsCommand struct {
	service *Service
	In      *api.GetTeamsRequest
	Out     *api.Teams
}

func NewGetTeamsCommand(service *Service, in *api.GetTeamsRequest) *GetTeamsCommand {
	return &GetTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	cursor, err := c.service.db.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	c.Out, err = toTeams(ctx, cursor, c.In.Page, c.In.Max)
	if err != nil {
		return err
	}
	return nil
}
