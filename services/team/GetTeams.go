package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
)

type GetTeamsCommand struct {
	service *TeamService
	In      *schema.GetTeamsRequest
	Out     *schema.Teams
}

func NewGetTeamsCommand(service *TeamService, in *schema.GetTeamsRequest) *GetTeamsCommand {
	return &GetTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	cursor, err := c.service.db.Aggregate(ctx, c.service.pipeline)
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
